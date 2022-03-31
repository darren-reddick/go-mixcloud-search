package mixcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

const (
	DefaultConcurrency int    = 5
	DefaultPageLimit   int    = 100
	DefaultHost        string = "api.mixcloud.com"
)

type config struct {
	Concurrency int
	PageLimit   int
	Limit       int // the limit for individual mixes returned in total
}

func NewConfig() config {
	return config{
		DefaultConcurrency,
		DefaultPageLimit,
		0,
	}
}

type Paging struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Response struct {
	Data   []Mix `json:"data,omitempty"`
	Paging `json:"paging,omitempty"`
}

type Mix struct {
	Key  string `json:"key,omitempty"`
	Url  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type ClientResponse struct {
	Mixes []Mix
	Next  *url.URL
}

const (
	MixSearch = iota
	ListenSearch
)

type ClientIface interface {
	Get(s string) (resp *http.Response, err error)
}

type Search struct {
	Term string
	Filter
	Client ClientIface
	Url    url.URL
	Store
	config
	logger *zap.Logger
}

type invalidSearchTermError struct {
	term string
	msg  string
}

func (i *invalidSearchTermError) Error() string {
	return fmt.Sprintf("%s: %s", i.term, i.msg)
}

func validateSearchTerm(s string) error {
	if s == "" {
		return &invalidSearchTermError{s, "Search term is invalid"}
	}
	return nil
}

func NewMixSearch(s string, filter Filter, client ClientIface, store Store, logger *zap.Logger) (Search, error) {
	err := validateSearchTerm(s)
	if err != nil {
		return Search{}, err
	}
	u := url.URL{
		Scheme:   "https",
		Host:     DefaultHost,
		Path:     "/search/",
		RawQuery: "type=cloudcast",
	}

	q := u.Query()
	q.Add("q", s)
	u.RawQuery = q.Encode()
	return Search{
		s,
		filter,
		client,
		u,
		store,
		NewConfig(),
		logger,
	}, nil
}

func NewHistorySearch(user string, filter Filter, client ClientIface, store Store, logger *zap.Logger) (Search, error) {

	u := url.URL{
		Scheme:   "https",
		Host:     DefaultHost,
		Path:     fmt.Sprintf("/%s/listens/", user),
		RawQuery: "",
	}

	return Search{
		"",
		filter,
		client,
		u,
		store,
		NewConfig(),
		logger,
	}, nil
}

func (a *Search) Get(offset int, limit int) (bool, error) {

	more := false

	u := a.Url
	q := u.Query()
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	a.logger.Debug(u.String())

	resp, err := a.Client.Get(u.String())
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	r := Response{}
	json.Unmarshal(b, &r)

	for _, mix := range a.Filter.Filter(r.Data) {
		err = a.Put(mix)
		if err != nil {
			return false, err
		}
	}

	if r.Paging.Next != "" {
		more = true
	}

	return more, nil
}

func (a *Search) GetAllParallel() error {
	offset := 0
	complete := false
	var wg sync.WaitGroup
	completeChan := make(chan bool, a.config.Concurrency)

	for !complete {
		for i := 1; i <= a.config.Concurrency; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				var more bool
				o := offset + ((i - 1) * a.config.PageLimit)
				a.logger.Debug(fmt.Sprintf("Thread %d fetching %d", i, o))

				more, _ = a.Get(o, a.config.PageLimit)

				if !more {
					completeChan <- true
				}

			}(i)
		}

		wg.Wait()
		select {
		case complete = <-completeChan:
			a.logger.Debug("complete signal received")
		default:
		}

		offset += (a.config.Concurrency * a.config.PageLimit)
	}

	return nil

}

func (a *Search) WriteJsonToFile() error {
	var data []byte
	data, err := json.MarshalIndent(&a.Data, "", "    ")
	if err != nil {

		return err

	}

	err = ioutil.WriteFile("test.json", data, 0644)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return nil
}

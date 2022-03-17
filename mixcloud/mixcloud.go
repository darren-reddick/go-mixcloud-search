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
)

const (
	DefaultConcurrency int = 5
	DefaultLimit       int = 100
)

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

var config = struct {
	Host  string
	Path  string
	Query string
}{
	"api.mixcloud.com",
	"/search/",
	"type=cloudcast&limit=100",
}

type ClientIface interface {
	Get(s string) (resp *http.Response, err error)
}

type Search struct {
	Term string
	Filter
	Client ClientIface
	Url    url.URL
	Store
}

func NewSearch(s string, filter Filter, client ClientIface, store Store) Search {
	u := url.URL{
		Scheme:   "https",
		Host:     config.Host,
		Path:     config.Path,
		RawQuery: config.Query,
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
	}
}

func (a *Search) Get(offset int) (bool, error) {

	more := false

	u := a.Url
	q := u.Query()
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	resp, err := a.Client.Get(u.String())
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	r := Response{}
	json.Unmarshal(b, &r)

	for _, mix := range a.Filter.Filter(r.Data) {
		a.Put(mix)
	}

	if r.Paging.Next != "" {
		more = true
	}

	return more, nil
}

func (a *Search) GetAllSync() error {
	offset := 0
	more := true
	var err error

	for more == true {
		fmt.Printf("Fetching offset %d\n", offset)
		more, err = a.Get(offset)
		if err != nil {
			return err
		}
		offset += 100
	}

	return nil

}

func (a *Search) GetAllAsync() error {
	offset := 0
	complete := false
	var err error
	var wg sync.WaitGroup
	completeChan := make(chan bool, DefaultConcurrency)

	for complete == false {
		for i := 1; i <= DefaultConcurrency; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				var more bool
				o := offset + ((i - 1) * DefaultLimit)
				fmt.Printf("Fetching %d\n", o)

				more, err = a.Get(o)

				if !more {
					completeChan <- true
				}

			}(i)
		}

		wg.Wait()
		fmt.Println("Done waiting")
		select {
		case complete = <-completeChan:
			fmt.Println("complete signal received")
		default:
		}

		if err != nil {
			return err
		}
		offset += (DefaultConcurrency * DefaultLimit)
	}

	return nil

}

func (a *Search) WriteJsonToFile() error {
	data := []byte{}
	data, err := json.MarshalIndent(&a.Data, "", "    ")
	if err != nil {
		fmt.Println("Write to file")

		return err

	}

	err = ioutil.WriteFile("test.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

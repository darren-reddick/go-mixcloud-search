package mixcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	"type=cloudcast&limit=100&offset=1000",
}

type HttpIface interface {
	Get(s string) (resp *http.Response, err error)
}

type Filter struct {
	Include []string
	Exclude []string
}

type Mixcloud struct {
	Search string
	Filter
	Http HttpIface
	Url  url.URL
	Store
}

func NewMixcloud(s string, filter Filter, http HttpIface, store Store) Mixcloud {
	u := url.URL{
		Scheme:   "https",
		Host:     config.Host,
		Path:     config.Path,
		RawQuery: config.Query,
	}

	q := u.Query()
	q.Add("q", s)
	u.RawQuery = q.Encode()
	return Mixcloud{
		s,
		filter,
		http,
		u,
		store,
	}
}

func (a *Mixcloud) Get(offset int) (bool, error) {

	more := false

	u := a.Url
	q := u.Query()
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	resp, err := a.Http.Get(u.String())
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	r := Response{}
	json.Unmarshal(b, &r)

	for _, mix := range r.Data {
		a.Put(mix)
	}

	if r.Paging.Next != "" {
		more = true
	}

	return more, nil
}

func (a *Mixcloud) GetAll() error {
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
		time.Sleep(3 * time.Second)
	}

	return nil

}

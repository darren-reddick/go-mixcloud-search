package mixcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/darren-reddick/go-mixcloud-search/schema"
	"github.com/darren-reddick/go-mixcloud-search/store"
)

type Paging struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Response struct {
	Data   []schema.Mix `json:"data,omitempty"`
	Paging `json:"paging,omitempty"`
}

type ClientResponse struct {
	Mixes []schema.Mix
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
	store.Store
}

func NewMixcloud(s string, filter Filter, http HttpIface, store store.Store) Mixcloud {
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

func (a *Mixcloud) Get() (ClientResponse, error) {
	cr := ClientResponse{}
	resp, err := a.Http.Get(a.Url.String())
	if err != nil {
		return cr, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	r := Response{}
	json.Unmarshal(b, &r)

	for _, mix := range r.Data {
		a.Put(mix)
	}

	if r.Paging.Next != "" {
		cr.Next, err = url.Parse(r.Paging.Next)
		if err != nil {
			return cr, err
		}
	}

	fmt.Println(a.Data)

	return cr, nil
}

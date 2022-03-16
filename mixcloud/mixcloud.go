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

func (a *Mixcloud) GetAllSync() error {
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

func (a *Mixcloud) GetAllAsync() error {
	offset := 0
	complete := false
	var err error
	var wg sync.WaitGroup
	completeChan := make(chan bool, 5)

	for complete == false {
		for i := 1; i <= 5; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				var more bool
				o := offset + ((i - 1) * 100)
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
			fmt.Println("No signal received")
		}

		if err != nil {
			return err
		}
		offset += 500
		//return nil
	}

	return nil

}

func (a *Mixcloud) WriteJsonToFile() error {
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

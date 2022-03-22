package mixcloud

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

type MockClient struct {
	resp *http.Response
}

func NewMockClient(body string) *MockClient {
	return &MockClient{
		&http.Response{
			Status:           "",
			StatusCode:       0,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           map[string][]string{},
			Body:             io.NopCloser(strings.NewReader(body)),
			ContentLength:    0,
			TransferEncoding: []string{},
			Close:            false,
			Uncompressed:     false,
			Trailer:          map[string][]string{},
			Request:          &http.Request{},
			TLS:              &tls.ConnectionState{},
		},
	}
}

func (m *MockClient) Get(s string) (resp *http.Response, err error) {
	return m.resp, nil
}

type MockPagingClient struct {
	Data []Response
	resp *http.Response
	*sync.RWMutex
}

func NewMockPagingClient(pageCount int, pageLength int) MockPagingClient {
	return MockPagingClient{
		GeneratePages(pageCount, pageLength),
		&http.Response{
			Status:           "",
			StatusCode:       0,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           map[string][]string{},
			Body:             io.NopCloser(strings.NewReader("")),
			ContentLength:    0,
			TransferEncoding: []string{},
			Close:            false,
			Uncompressed:     false,
			Trailer:          map[string][]string{},
			Request:          &http.Request{},
			TLS:              &tls.ConnectionState{},
		},
		&sync.RWMutex{},
	}

}

func (m *MockPagingClient) Get(s string) (resp *http.Response, err error) {
	m.Lock()
	defer m.Unlock()
	j, _ := json.MarshalIndent(m.Data[0], "", "    ")
	m.resp.Body = io.NopCloser(strings.NewReader(string(j)))

	if len(m.Data) > 0 {
		m.Data = m.Data[1:]
	}

	return m.resp, nil

}

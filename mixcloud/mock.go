package mixcloud

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
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

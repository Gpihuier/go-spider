package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	req    *http.Request
	Err    error
	client http.Client
	resp   *http.Response
}

func ParseUrl(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	return fmt.Sprintf("http://%s", url)
}

func NewRequest(method, url string) *Request {
	req, err := http.NewRequest(method, ParseUrl(url), nil)
	return &Request{
		req:    req,
		Err:    err,
		client: http.Client{},
	}
}

func Get(url string) *Request {
	return NewRequest(http.MethodGet, url)
}

func Post(url string) *Request {
	return NewRequest(http.MethodPost, url)
}

func (s *Request) AddCookie(c *http.Cookie) *Request {
	s.req.AddCookie(c)
	return s
}

func (s *Request) AddHeader(key, value string) *Request {
	s.req.Header.Add(key, value)
	return s
}

func (s *Request) AddHeaders(v map[string]string) *Request {
	for k, v := range v {
		s.AddHeader(k, v)
	}
	return s
}

func (s *Request) SetUA(ua string) *Request {
	s.AddHeader("User-Agent", ua)
	return s
}

func (s *Request) SetRandomUA() *Request {
	s.AddHeader("User-Agent", RandomUa())
	return s
}

func (s *Request) AddParam(k, v string) *Request {
	if len(s.req.URL.RawQuery) > 0 {
		s.req.URL.RawQuery += "&"
	}
	s.req.URL.RawQuery += url.QueryEscape(k) + "=" + url.QueryEscape(v)
	return s
}

func (s *Request) AddParams(v map[string]string) *Request {
	for k, v := range v {
		s.AddParam(k, v)
	}
	return s
}

func (s *Request) SetBody(b io.Reader) *Request {

	rc, ok := b.(io.ReadCloser)
	if !ok && b != nil {
		rc = io.NopCloser(b)
	}
	s.req.Body = rc

	switch v := b.(type) {
	case *bytes.Buffer:
		s.req.ContentLength = int64(v.Len())
		buf := v.Bytes()
		s.req.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return io.NopCloser(r), nil
		}
	case *bytes.Reader:
		s.req.ContentLength = int64(v.Len())
		snapshot := *v
		s.req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	case *strings.Reader:
		s.req.ContentLength = int64(v.Len())
		snapshot := *v
		s.req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	default:
	}
	return s
}

func (s *Request) SetRawBody(b []byte) *Request {
	s.SetBody(bytes.NewReader(b))
	return s
}

func (s *Request) SetFormBody(v map[string]string) *Request {
	var u url.URL
	q := u.Query()
	for k, v := range v {
		q.Add(k, v)
	}
	s.SetRawBody([]byte(q.Encode()))
	s.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	return s
}

func (s *Request) SetJsonBody(v interface{}) *Request {
	body, err := json.Marshal(v)
	s.SetRawBody(body)
	s.Err = err
	s.AddHeader("Content-Type", "application/json")
	return s
}

func (s *Request) Do() *http.Response {
	s.resp, s.Err = s.client.Do(s.req)
	return s.resp
}

func (s *Request) String() string {
	return s.req.URL.String()
}

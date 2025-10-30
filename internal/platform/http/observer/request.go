package observer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// io.TeeReader

type RequestObserver struct {
	*http.Request

	body *bytes.Buffer

	routerFn func(*http.Request) string
}

// NewRequestObserver creates a new *RequestObserver instance.
//
// Parameters:
//   - r *http.Request: request;
//   - readBody bool: indicates whether the request body should be read.
func NewRequestObserver(r *http.Request, readBody bool, routerFn func(*http.Request) string) *RequestObserver {
	obs := &RequestObserver{
		Request:  r,
		body:     &bytes.Buffer{},
		routerFn: routerFn,
	}

	if readBody {
		if r.Body != nil {
			_, _ = obs.body.ReadFrom(obs.Body)

			obs.Body = io.NopCloser(bytes.NewReader(obs.body.Bytes()))
		}
	}

	return obs
}

func (r *RequestObserver) GetBodyString() string {
	return r.body.String()
}

// content_length
func (r *RequestObserver) GetBodySize() int64 {
	if r.ContentLength >= 0 {
		return r.ContentLength
	}

	counter := &countingReader{r: r.Body}
	r.Body = counter
	return -1
}

func (r *RequestObserver) GetMethod() string {
	return r.Method
}

func (r *RequestObserver) GetURLPath() string {
	return r.URL.Path
}

func (r *RequestObserver) GetURLQuery() string {
	return r.URL.RawQuery
}

func (r *RequestObserver) GetProtocol() string {
	return r.Proto
}

// шаблон маршрута
func (r *RequestObserver) GetRoute() string {
	//return chi.RouteContext(r.Context()).RoutePattern()
	return r.routerFn(r.Request)
}

func (r *RequestObserver) GetURI() string {
	return fmt.Sprintf("%s %s%s %s",
		r.Method,
		r.URL.Path,
		func() string {
			if r.URL.RawQuery != "" {
				return "?" + r.URL.RawQuery
			}

			return ""
		}(),
		r.Proto,
	)
}

type countingReader struct {
	r     io.ReadCloser
	count int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.count += int64(n)

	return n, err
}

func (c *countingReader) Close() error {
	return c.r.Close()
}

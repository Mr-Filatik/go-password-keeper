// Package observer provides functionality for wrapper structures that output information.
package observer

import (
	"bytes"
	"io"
	"net/http"
)

// RouteFunc describes the type of function for getting a route.
type RouteFunc func(*http.Request) string

func defaultRouteFunc() RouteFunc {
	return func(_ *http.Request) string {
		return "unknown"
	}
}

// RequestObserver - a structure for obtaining information about a request.
// Intercepts the content to obtain information about it and passes it on.
type RequestObserver struct {
	*http.Request

	bodyBuf  *bytes.Buffer
	bodySize int64
	routeFn  RouteFunc
}

const maxBufferedRequestBody = 64 << 10

// NewRequestObserver creates a new *RequestObserver instance.
//
// Parameters:
//   - r *http.Request: request;
//   - readBody bool: indicates whether the request body should be read.
func NewRequestObserver(r *http.Request, readBody bool, routeFn RouteFunc) *RequestObserver {
	obs := &RequestObserver{
		Request:  r,
		bodyBuf:  nil,
		bodySize: -1,
		routeFn:  routeFn,
	}

	if obs.routeFn == nil {
		obs.routeFn = defaultRouteFunc()
	}

	if r.ContentLength >= 0 {
		obs.bodySize = r.ContentLength
	}

	if readBody && r.Body == nil {
		return obs
	}

	defer func() {
		_ = r.Body.Close()
	}()

	var buf bytes.Buffer

	n, _ := io.CopyN(&buf, r.Body, maxBufferedRequestBody+1)

	if obs.bodySize < 0 {
		if n <= maxBufferedRequestBody {
			obs.bodySize = n
		} else {
			obs.bodySize = int64(maxBufferedRequestBody) + 1
		}
	}

	if buf.Len() > maxBufferedRequestBody {
		buf.Truncate(maxBufferedRequestBody)
	}

	obs.bodyBuf = &buf

	r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

	return obs
}

// GetBodyString returns the request body as a string.
func (r *RequestObserver) GetBodyString() string {
	if r.bodyBuf == nil {
		return ""
	}

	return r.bodyBuf.String()
}

// GetBodySize returns the size of the request body.
func (r *RequestObserver) GetBodySize() int64 {
	return r.bodySize
}

// GetMethod returns the request method.
func (r *RequestObserver) GetMethod() string {
	return r.Method
}

// GetURLPath returns the request path.
func (r *RequestObserver) GetURLPath() string {
	return r.URL.Path
}

// GetURLQuery returns the query parameters.
func (r *RequestObserver) GetURLQuery() string {
	return r.URL.RawQuery
}

// GetProtocol returns the request protocol.
func (r *RequestObserver) GetProtocol() string {
	return r.Proto
}

// GetHeader returns the value for a request header by key.
//
// Parameters:
//   - hdrName string: header name.
func (r *RequestObserver) GetHeader(hdrName string) string {
	return r.Header.Get(hdrName)
}

// GetRoute returns a route template.
// Unlike GetURLPath(), it does not include mutable data in the request.
func (r *RequestObserver) GetRoute() string {
	if r.routeFn == nil {
		return "unknown"
	}

	route := r.routeFn(r.Request)
	if route == "" {
		return "/"
	}

	return route
}

// GetURI returns the URI for the request.
func (r *RequestObserver) GetURI() string {
	return r.Method + " " + r.URL.RequestURI() + " " + r.Proto
}

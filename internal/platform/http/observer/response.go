// Package observer provides functionality for wrapper structures that output information.
package observer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// ResponseObserver - a structure for obtaining information about a response.
// Intercepts the content to obtain information about it and passes it on.
type ResponseObserver struct {
	http.ResponseWriter

	status  int
	size    int64
	bodyBuf *bytes.Buffer

	writer io.Writer
}

// NewResponseObserver creates a new *ResponseObserver instance.
//
// Parameters:
//   - w http.ResponseWriter: response writer;
//   - readBody bool: indicates whether the request body should be read.
func NewResponseObserver(w http.ResponseWriter, readBody bool) *ResponseObserver {
	obs := &ResponseObserver{
		ResponseWriter: w,
		status:         0,
		size:           0,
		bodyBuf:        nil,
		writer:         w,
	}

	if readBody {
		obs.bodyBuf = &bytes.Buffer{}
		obs.writer = io.MultiWriter(w, obs.bodyBuf)
	}

	return obs
}

// Header returns a map of response headers.
//
// Implements the http.ResponseWriter interface.
func (r *ResponseObserver) Header() http.Header {
	return r.ResponseWriter.Header()
}

// WriteHeader writes the response code.
//
// Implements the http.ResponseWriter interface.
func (r *ResponseObserver) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Write writes the contents of the response body.
//
// Implements the http.ResponseWriter interface.
func (r *ResponseObserver) Write(bytes []byte) (int, error) {
	num, err := r.writer.Write(bytes)
	r.size += int64(num)

	if err != nil {
		return num, fmt.Errorf("write: %w", err)
	}

	return num, nil
}

// GetStatus returns the response status.
func (r *ResponseObserver) GetStatus() int {
	return r.status
}

// GetBodyString returns the response body as a string.
func (r *ResponseObserver) GetBodyString() string {
	if r.bodyBuf == nil {
		return ""
	}

	return r.bodyBuf.String()
}

// GetBodySize returns the size of the response body.
func (r *ResponseObserver) GetBodySize() int64 {
	return r.size
}

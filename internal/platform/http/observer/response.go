package observer

import (
	"bytes"
	"io"
	"net/http"
)

// TeeReader

type ResponseObserver struct {
	http.ResponseWriter

	status int
	size   int64
	body   *bytes.Buffer
	writer io.Writer
}

func NewResponseObserver(w http.ResponseWriter, captureBody bool) *ResponseObserver {
	obs := &ResponseObserver{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
	if captureBody {
		obs.body = &bytes.Buffer{}
		obs.writer = io.MultiWriter(w, obs.body)
	} else {
		obs.writer = w
	}
	return obs
}

func (r *ResponseObserver) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseObserver) Write(b []byte) (int, error) {
	n, err := r.writer.Write(b)
	r.size += int64(n)

	return n, err
}

func (r *ResponseObserver) GetStatus() int {
	return r.status
}

// в байтах выдаёт значение без копирования
func (r *ResponseObserver) GetBodyBytes() []byte {
	return r.body.Bytes()
}

// в utf8 выдаёт значение + копирование
// для больших ответов лучше отключить
func (r *ResponseObserver) GetBodyString() string {
	return r.body.String()
}

func (r *ResponseObserver) GetBodySize() int64 {
	return r.size
}

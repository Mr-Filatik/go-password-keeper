// Package mocks contains stubs for other tests.
package mocks

import (
	"bytes"
	"fmt"
	"sync"
)

// MockWriter implementation of oi.Writer for tests.
type MockWriter struct {
	// mu - mutex.
	mu sync.Mutex

	// buffer - buffer for storing data.
	buffer *bytes.Buffer

	// datas - data in the form of a []byte slice.
	datas [][]byte

	// countUnreadedDatas - the number of unread data.
	countUnreadedDatas int
}

// NewMockWriter initializes and creates a new *MockWriter instance.
func NewMockWriter() *MockWriter {
	return &MockWriter{
		mu:                 sync.Mutex{},
		buffer:             &bytes.Buffer{},
		datas:              make([][]byte, 0),
		countUnreadedDatas: 0,
	}
}

// Write writes data to the buffer and slice.
//
// Implements the io.Writer interface.
func (w *MockWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	num, err := w.buffer.Write(data)
	if err != nil {
		return num, fmt.Errorf("buffer write: %w", err)
	}

	w.datas = append(w.datas, bytes.TrimSpace(data))

	w.countUnreadedDatas++

	return num, nil
}

// MarkDataAsRead marks all unread data as read.
func (w *MockWriter) MarkDataAsRead() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.countUnreadedDatas = 0
}

// GetUnreadedData returns the latest unread data and a flag indicating whether it exists.
func (w *MockWriter) GetUnreadedData() ([]byte, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.countUnreadedDatas == 0 {
		return nil, false
	}

	data := w.datas[len(w.datas)-w.countUnreadedDatas]

	w.countUnreadedDatas--

	return data, true
}

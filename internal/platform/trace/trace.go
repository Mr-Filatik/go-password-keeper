package trace

import (
	"crypto/rand"
	"fmt"
)

type Trace struct {
	traceID  string
	parentID string
	spanID   string
	remote   bool
}

// New сосздаёт новый экземпляр для отслеживания трейсинга.
//
// Создание используется только в middleware, поэтому значение параметра remote выставляется в true.
func New(traceID, spanID string) Trace {
	return Trace{
		traceID:  traceID,
		parentID: spanID,
		spanID:   generateSpanID(),
		remote:   true,
	}
}

func NewEmpty() Trace {
	return Trace{
		traceID:  generateTraceID(),
		parentID: "",
		spanID:   generateSpanID(),
		remote:   false,
	}
}

func (t Trace) Next() Trace {
	return Trace{
		traceID:  t.traceID,
		parentID: t.spanID,
		spanID:   generateSpanID(),
		remote:   false,
	}
}

func (t Trace) TraceID() string {
	return t.traceID
}

func (t Trace) ParentID() string {
	return t.parentID
}

func (t Trace) SpanID() string {
	return t.spanID
}

func (t Trace) Remote() bool {
	return t.remote
}

func generateTraceID() string {
	const lenTraceID = 16

	bytes := make([]byte, lenTraceID)

	_, err := rand.Read(bytes)
	if err != nil {
		return "0000000000000000"
	}

	return fmt.Sprintf("%032x", bytes)
}

func generateSpanID() string {
	const lenSpanID = 8

	bytes := make([]byte, lenSpanID)

	_, err := rand.Read(bytes)
	if err != nil {
		return "0000000000000000"
	}

	return fmt.Sprintf("%016x", bytes)
}

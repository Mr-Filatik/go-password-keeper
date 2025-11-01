// Package logging provides logging functionality.
package logging

// Main fields for logging.
const (
	FieldBaseTimestamp = "timestamp"
	FieldBaseLevel     = "level"
	FieldBaseMessage   = "message"
	FieldBaseCaller    = "caller"
)

// Fields used for tracing.
const (
	FieldRequestID = "request.id"
	FieldSpanID    = "span.id"
	FieldTraceID   = "trace.id"
)

// HTTP-related fields for logging.
const (
	FieldHTTPMethod       = "http.method"
	FieldHTTPPath         = "http.path"
	FieldHTTPQuery        = "http.query"
	FieldHTTPRoute        = "http.route"
	FieldHTTPStatusCode   = "http.status_code"
	FieldHTTPRequestSize  = "http.request_size"
	FieldHTTPResponceSize = "http.response_size"
	FieldHTTPDurationMs   = "http.duration_ms"
	FieldClientIP         = "client.ip"
)

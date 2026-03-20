// Package logging provides logging functionality.
package logging

// Main fields for logging.
const (
	FieldBaseTimestamp      = "timestamp" // Logging time.
	FieldBaseTimestampShort = "ts"        // Logging time.
	FieldBaseLevel          = "level"     // Log level.
	FieldBaseLevelShort     = "lvl"       // Log level.
	FieldBaseMessage        = "message"   // Message indicating the reason for the log, error, etc.
	FieldBaseMessageShort   = "msg"       // Message indicating the reason for the log, error, etc.
	FieldBaseError          = "error"     // Error.
	FieldBaseErrorShort     = "err"       // Error.

	FieldBaseCaller     = "caller" // Code call location.
	FieldBaseStackTrace = "stacktrace"

	// Данные внутри приложения.
	FieldBaseData = "data" // Additional data not included in the message text.
	// Данные, которые передаются или принимаются откуда-то.
	FieldBasePayload = "payload" // Additional data not included in the message text.
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

	FieldClientIP = "client.ip"
)

//nolint:godox
// TODO: Convert to a structure to standardize log fields.
// log.With([]log.Field{
// 		log.Stack("stacktrace"),
// 		log.String("request_id", middleware.GetRequestIDFromContext(ctx)),
// 		log.Err(err.(error)),
// }...).Error("blet gql handler panic with stacktrace")

// Package logging provides logging functionality.
package log

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

	FieldBaseCaller    = "caller" // Code call location.
	FieldBaseCallStack = "callstack"

	// FieldBaseStackTrace is a special field for displaying the full call path.
	//
	// It should contain the full chain of function calls from main() to the current function.
	// This is quite resource-intensive; it is recommended to use this field only for panics.
	FieldBaseStackTrace = "stacktrace"

	// FieldBaseData field is a special field for additional data not included in the message text.
	//
	// This field is recommended for data used within the application.
	FieldBaseData = "data"

	// FieldBasePayload field is a special field for additional data not included in the message text.
	//
	// This field is recommended for use with data that is transmitted or received from somewhere external.
	FieldBasePayload = "payload"
)

const (
	FieldProject   = "project"
	FieldApp       = "app"
	FieldComponent = "component"
)

// Fields used for tracing.
const (
	FieldRequestID = "request.id"

	FieldTraceID  = "trace.id"
	FieldSpanID   = "span.id"
	FieldParentID = "parent.id"
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

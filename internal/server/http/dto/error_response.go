// Package dto contains the entities that the server works with.
package dto

// ErrorResponse describes the structure of the response when an error occurs.
//
// Used when the server response has an error code of 400 or higher.
type ErrorResponse struct {
	// StatusCode - HTTP response status code.
	//
	// It is a required parameter.
	StatusCode int `json:"statusCode" validate:"required,min=400,max=599"`

	// Error - error.
	//
	// It is a required parameter.
	Error Error `json:"error" validate:"required"`
} // @name ErrorResponse

// Error describes detailed information about the error.
type Error struct {
	// ErrorType - machine-readable error type.
	//
	// It is a required parameter.
	Type ErrorType `json:"type" validate:"required"`

	// ID - unique error identifier.
	//
	// It is a required parameter.
	ID string `json:"id" validate:"required"`

	// Message - main error message.
	//
	// It is a required parameter.
	Message string `json:"message" validate:"required"`

	// Details - additional details about the error.
	//
	// This is an optional parameter.
	Details []string `json:"details,omitempty" validate:"omitempty"`
} // @name Error

// ErrorType - machine-readable error type.
type ErrorType string // @name ErrorType

const (
	// RequestInvalidFormat - invalid request model format.
	RequestInvalidFormat ErrorType = "REQUEST_INVALID_FORMAT"

	// RequestValidationError - query model validation errors.
	RequestValidationError ErrorType = "REQUEST_VALIDATION_ERROR"
)

// Package handler contains all server handlers.
package handler

import (
	"encoding/json"
	"net/http"

	//"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/platform/validator"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/dto"
)

const megaByte = 1024 * 1024 // 1MB

func getRequest[T any](
	w http.ResponseWriter,
	r *http.Request,
	validator validator.IValidator,
) (*T, bool) { // option.WithSkipValidate, option.WithValidator
	//ctx := r.Context()
	// Checking Content-Type
	r.Body = http.MaxBytesReader(w, r.Body, megaByte) // Limiting body size (DOS protection)

	defer func() {
		err := r.Body.Close()
		if err != nil {
			//logging.Error(ctx, "Body close failed", err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req T

	decodeErr := decoder.Decode(&req)
	if decodeErr != nil {
		sendError(w, r,
			http.StatusBadRequest,
			dto.RequestInvalidFormat, "Invalid model when requesting",
			nil,
		)

		return &req, false
	}

	problems, err := validator.Validate(req)
	if err != nil {
		sendError(w, r,
			http.StatusBadRequest,
			dto.RequestValidationError, "Validation error(s)",
			problems,
		)

		return &req, false
	}

	return &req, true
}

func sendError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	errorType dto.ErrorType,
	message string,
	details []string,
) {
	//ctx := r.Context()

	response := &dto.ErrorResponse{
		StatusCode: status,
		Error: dto.Error{
			ID:      "none",
			Type:    errorType,
			Message: message,
			Details: details,
		},
	}

	jsonResponse, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		//logging.Error(ctx, "Marshal response failed", marshalErr)

		http.Error(w,
			"An unexpected internal server error occurred", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	_, writeErr := w.Write(jsonResponse)
	if writeErr != nil {
		//logging.Error(ctx, "Failed to write error response", writeErr)
	}
}

func sendSuccess(
	w http.ResponseWriter,
	r *http.Request,
	response any,
) {
	//ctx := r.Context()

	jsonResponse, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		//logging.Error(ctx, "Marshal response failed", marshalErr)

		http.Error(w,
			"An unexpected internal server error occurred", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	_, writeErr := w.Write(jsonResponse)
	if writeErr != nil {
		//logging.Error(ctx, "Failed to write error response", writeErr)
	}
}

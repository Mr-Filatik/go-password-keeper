// Package handler contains all server handlers.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/dto"
)

func getRequest[T any](
	w http.ResponseWriter,
	r *http.Request,
) (*T, bool) { // option.WithSkipValidate
	// 1. Проверяем Content-Type
	// 2. Ограничиваем размер тела (защита от DOS)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req T
	if err := decoder.Decode(&req); err != nil {
		// s.logger.Error("Decode error", err)
		//return nil, fmt.Errorf("%w: decode failed: %s", ErrRequestFormat, err.Error())
		sendError(w, r,
			http.StatusBadRequest,
			dto.RequestInvalidFormat, "Invalid model when requesting",
			nil,
		)
		return &req, false
	}

	validate := validator.New()

	err := validate.Struct(req)
	if err != nil {
		// s.logger.Error("Validate error", err)
		var details []string
		// Проверяем тип ошибки
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				// fmt.Printf("Поле: %s\n", e.Field())
				// fmt.Printf("Тэг: %s\n", e.Tag())
				// fmt.Printf("Значение: %v\n", e.Value())
				// fmt.Printf("Параметр: %s\n", e.Param())
				// fmt.Printf("Тип: %v\n", e.Type())
				// fmt.Printf("Ошибка: %s\n", e.Error())
				// fmt.Println("---")
				details = append(details, getValidationMessage(fieldError))
			}
		}
		sendError(w, r,
			http.StatusBadRequest,
			dto.RequestValidationError, "Validation error(s)",
			details,
		)
		return &req, false
		//return nil, fmt.Errorf("%w: validation failed: %s", ErrRequestValidation, err.Error())
	}

	return &req, true
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", e.Field())

	case "min":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at least %s characters long",
				e.Field(), e.Param())
		}

		return fmt.Sprintf("Field '%s' must be >= %s", e.Field(), e.Param())

	case "max":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at most %s characters long",
				e.Field(), e.Param())
		}

		return fmt.Sprintf("Field '%s' must be <= %s", e.Field(), e.Param())

	case "gte":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at least %s characters long",
				e.Field(), e.Param())
		}

		return fmt.Sprintf("Field '%s' must be >= %s", e.Field(), e.Param())

	case "lte":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at most %s characters long",
				e.Field(), e.Param())
		}

		return fmt.Sprintf("Field '%s' must be <= %s", e.Field(), e.Param())

	case "gt":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be more than %s characters long",
				e.Field(), e.Param())
		}

		return fmt.Sprintf("Field '%s' must be > %s", e.Field(), e.Param())

	case "lt":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be less than %s characters long",
				e.Field(), e.Param())
		}

		return fmt.Sprintf("Field '%s' must be < %s", e.Field(), e.Param())

	case "len":
		return fmt.Sprintf("Field '%s' must be exactly %s characters long", e.Field(), e.Param())

	case "numeric":
		return fmt.Sprintf("Field '%s' must contain only numbers", e.Field())

	case "alpha":
		return fmt.Sprintf("Field '%s' must contain only letters", e.Field())

	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only letters and numbers", e.Field())

	case "email":
		return fmt.Sprintf("Field '%s' must be in email format", e.Field())

	case "uuid":
		return fmt.Sprintf("Field '%s' must be in UUID format", e.Field())

	case "url":
		return fmt.Sprintf("Field '%s' must be a valid URL", e.Field())

	case "ip":
		return fmt.Sprintf("Field '%s' must be a valid IP address", e.Field())

	// case "oneof": return fmt.Sprintf("Field '%s' must be one of: %s", e.Field(), e.Param())

	case "datetime":
		return fmt.Sprintf("Field '%s' must be in format: %s", e.Field(), e.Param())

	default:
		return e.Error() // unknown tag
	}
}

func sendError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	errorType dto.ErrorType,
	message string,
	details []string,
) {
	ctx := r.Context()

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
		// Если не смогли даже ошибку сформировать - пишем plain text
		// s.logger.Error("Failed to marshal error response", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	_, writeErr := w.Write(jsonResponse)
	if writeErr != nil {
		logging.Error(ctx, "Failed to write error response", writeErr)
	}
}

func sendSuccess(
	w http.ResponseWriter,
	r *http.Request,
	response any,
) {
	ctx := r.Context()

	jsonResponse, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		// Если не смогли даже ошибку сформировать - пишем plain text
		// s.logger.Error("Failed to marshal error response", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	_, writeErr := w.Write(jsonResponse)
	if writeErr != nil {
		logging.Error(ctx, "Failed to write error response", writeErr)
	}
}

package handler

import (
	"context"
	"errors"
	"net/http"

	logctx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/log"
	tracectx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx/trace"
	"github.com/mr-filatik/go-password-keeper/internal/platform/types"
	"github.com/mr-filatik/go-password-keeper/internal/platform/validator"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/dto"
)

// Test describes a test server handler for development.
//
//	@Summary		Test handler
//	@Description	Test handler for developing
//	@Tags			developing
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.TestRequest		true	"Request"
//	@Success		200		{object}	dto.TestResponse	"Successful response"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format or validation error"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/test [post]
func Test(validator validator.IValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, err := dto.SecretTypeFromString("aaa")
		if err != nil {
			if errors.Is(err, types.ErrUnexpectedValue) {
				logctx.Error(ctx, "ErrUnexpectedValue", err)
			}

			SubProc(ctx, err)
		}

		req, ok := getRequest[dto.TestRequest](w, r, validator) // refactor
		if !ok {
			return
		}

		logctx.Info(ctx, "Test")

		req.Number++
		req.Message += " new"

		sendSuccess(w, r, &dto.TestResponse{
			Number:  req.Number,
			Message: req.Message,
			Mes:     req.Mes,
		})
	}
}

func SubProc(ctx context.Context, err error) {
	ctx = tracectx.Next(ctx)

	if errors.Is(err, types.ErrUnexpectedValueInEnum) {
		logctx.Error(ctx, "ErrUnexpectedValueInEnum", err)
	}

	if errors.Is(err, dto.ErrUnexpectedValueInEnumSecretType) {
		logctx.Error(ctx, "ErrUnexpectedValueInEnumSecretType", err)
	}
}

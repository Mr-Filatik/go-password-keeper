package logging

import (
	"fmt"
)

type FieldOption func(args []any) []any

func applyOptions(options ...FieldOption) []any {
	tempMap := make(map[any]struct{}, len(options))

	fields := make([]any, 0, len(options)*2+2) // Error + CallStack

	for _, applyOpt := range options {
		tempFields := applyOpt(fields)

		for i := 0; i < len(tempFields); i += 2 {
			_, ok := tempMap[tempFields[i]]
			if !ok {
				tempMap[tempFields[i]] = struct{}{}
				fields = applyOpt(fields)
			}
		}
	}

	return fields
}

// Стандартные ошибки в Go (созданные через errors.New) не содержат стек. Чтобы он появился, ошибку нужно «создать» или «обернуть» специальной библиотекой при возникновении:
// Создание: errors.WithStack(err) или errors.Errorf("...) из ://github.com.
type stackTracer interface {
	StackTrace() any // errors.StackTrace
}

// Хотите, чтобы я показал, как сделать
// WithErrorField универсальным, чтобы он умел доставать StackTrace из ошибки?
// Но тогда len(options) нарушится
func WithErrorField(err error) FieldOption {
	return func(args []any) []any {
		if st, ok := err.(stackTracer); ok {
			return append(args,
				FieldBaseError, err.Error(),
				FieldBaseStackTrace, fmt.Sprintf("%+v", st.StackTrace()))
		}

		return append(args, FieldBaseError, err.Error())
	}
}

func WithDataField(data any) FieldOption {
	return func(args []any) []any {
		return append(args, FieldBaseData, data)
	}
}

type ISanitizer interface {
	Sanitize() any
}

func WithDataFieldSanitised(data any) FieldOption {
	return func(args []any) []any {
		sanitised, ok := data.(ISanitizer)
		if !ok {
			return append(args, FieldBaseData, data)
		}

		return append(args, FieldBaseData, sanitised.Sanitize())
	}
}

func WithCustomField(key string, data any) FieldOption {
	return func(args []any) []any {
		return append(args, key, data)
	}
}

package mask

import (
	"errors"
	"fmt"
)

// ErrUnexpectedType - unexpected type during type conversion.
var ErrUnexpectedType = errors.New("unexpected type")

// MFunc represents a general function for masking data.
type MFunc func(any) (any, error)

// WrapStrToStrFn wraps a string-to-string function into a general wrapping function.
func WrapStrToStrFn(strfn func(string) string) MFunc {
	return func(data any) (any, error) {
		val, ok := data.(string)
		if !ok {
			return "", fmt.Errorf("%w: expected string", ErrUnexpectedType)
		}

		return strfn(val), nil
	}
}

// MaskedPassword wraps the password masking function to the type used in the package.
func MaskedPassword() MFunc {
	return WrapStrToStrFn(Password)
}

// MaskedEmail wraps the email masking function to the type used in the package.
func MaskedEmail() MFunc {
	return WrapStrToStrFn(Email)
}

// MaskedPhone wraps the phone masking function to the type used in the package.
func MaskedPhone() MFunc {
	return WrapStrToStrFn(Phone)
}

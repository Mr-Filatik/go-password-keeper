package mask

import (
	"errors"
	"fmt"
)

// ErrUnexpectedType - unexpected type during type conversion.
var ErrUnexpectedType = errors.New("unexpected type")

// EditFunc represents a general function for masking data.
type EditFunc func(any) (any, error)

type DeleteFunc func(any) (bool, error)

// WrapStrToStrFn wraps a string-to-string function into a general wrapping function.
func WrapStrToStrFn(strfn func(string) string) EditFunc {
	return func(data any) (any, error) {
		val, ok := data.(string)
		if !ok {
			return "", fmt.Errorf("%w: expected string", ErrUnexpectedType)
		}

		return strfn(val), nil
	}
}

// Password wraps the password masking function to the type used in the package.
func Password() EditFunc {
	return WrapStrToStrFn(PasswordMask)
}

// Email wraps the email masking function to the type used in the package.
func Email() EditFunc {
	return WrapStrToStrFn(EmailMask)
}

// Phone wraps the phone masking function to the type used in the package.
func Phone() EditFunc {
	return WrapStrToStrFn(PhoneMask)
}

func DAny() DeleteFunc {
	return func(_ any) (bool, error) {
		return true, nil
	}
}

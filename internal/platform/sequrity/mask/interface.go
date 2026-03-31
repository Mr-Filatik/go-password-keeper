// Package mask contains logic for masking sensitive data.
package mask

// IMasker describes a common interface for all maskers.
type IMasker interface {
	// Mask masks data represented as a slice of bytes according to the rules
	// described in the IMaskable interface for a specific type.
	Mask(data []byte, rules IMaskable) (string, error)
}

// IMaskable describes an interface for types that need to be masked when written somewhere.
type IMaskable interface {
	// Rules returns the rules for masking fields for the type.
	Rules() []Rule

	// replace on Mask(data []byte) for testing
	// нужно сразу иметь возможность тестами покрыть ошибки маскирования
}

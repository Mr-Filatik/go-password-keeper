// Package mask contains logic for masking sensitive data.
package mask

// IMasker describes a common interface for all maskers.
type IMasker interface {
	// Mask masks data represented as a slice of bytes according to the rules
	// described in the IMaskable interface for a specific type.
	Mask(data []byte, rules IMaskable) (string, error) // m.b. error to bool

	// MaskBytes(data []byte, rules []Rule) ([]byte, error)
}

// IMaskable describes an interface for types that need to be masked when written somewhere.
//
// Interface functions immediately return masked data to allow testing. This allows us to eliminate
// unnecessary errors in the application's operation by covering this area with tests.
type IMaskable interface {
	// Rules returns the rules for masking fields for the type.
	Rules() []Rule // remove rules in MaskBytes()

	// Mask masks data in the model, returning a new structure with the masked data.
	// Mask() any

	// MaskBytes masks serialized data according to type-specific rules.
	// MaskBytes(data []byte) []byte
}

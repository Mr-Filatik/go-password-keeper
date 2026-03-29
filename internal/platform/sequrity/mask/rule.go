package mask

// MType describes the field masking type.
type MType string

const (
	// MTypeNone - do not mask the field.
	MTypeNone MType = "none"

	// MTypeEdit - replace the value with a masked one.
	MTypeEdit MType = "edit"

	// MTypeEditInSlice - replace the value in the array with a masked one.
	MTypeEditInSlice MType = "edit-in-slice"

	// MTypeRemove - remove value.
	MTypeRemove MType = "remove"

	// MTypeRemoveInSlice - delete a value in an array.
	MTypeRemoveInSlice MType = "remove-in-slice"
)

// Rule describes the rule by which the field will be masked.
type Rule struct {
	path     []string
	intPath  []string
	maskType MType
	maskFunc MFunc
}

// GetPath returns the path to the field that needs to be masked.
func (r Rule) GetPath() []string {
	return r.path
}

// GetIntPath function returns the path to the nested field that needs to be masked.
func (r Rule) GetIntPath() []string {
	return r.intPath
}

// GetMaskType returns the type of masking to apply to the field.
func (r Rule) GetMaskType() MType {
	return r.maskType
}

// GetMaskFunc returns the function that should be used to mask the field.
func (r Rule) GetMaskFunc() MFunc {
	return r.maskFunc
}

// RuleEdit creates a rule to mask a field with a value replacement.
//
// The path parameter must match the names of the JSON tags in the structure.
func RuleEdit(path []string, fnc MFunc) Rule {
	return Rule{
		path:     path,
		intPath:  []string{},
		maskType: MTypeEdit,
		maskFunc: fnc,
	}
}

// RuleEditInSlice creates a rule to mask a field inside an array with a value replacement.
//
// The path and intPath parameters must match the names of the JSON tags in the structure.
func RuleEditInSlice(path, intPath []string, fnc MFunc) Rule {
	return Rule{
		path:     path,
		intPath:  intPath,
		maskType: MTypeEditInSlice,
		maskFunc: fnc,
	}
}

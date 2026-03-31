package mask

// Rule describes the rule by which the field will be masked.
type Rule struct {
	paths      []string
	maskFunc   EditFunc
	deleteFunc DeleteFunc
}

// GetPath returns the path to the field that needs to be masked.
func (r Rule) GetPath() []string {
	return r.paths
}

// GetMaskFunc returns the function that should be used to mask the field.
func (r Rule) GetMaskFunc() EditFunc {
	return r.maskFunc
}

func (r Rule) GetDeleteFunc() DeleteFunc {
	return r.deleteFunc
}

// RuleEdit creates a rule to mask a field with a value replacement.
//
// The path parameter must match the names of the JSON tags in the structure.
func RuleEdit(fnc EditFunc, paths ...string) Rule {
	return Rule{
		paths:      paths,
		maskFunc:   fnc,
		deleteFunc: nil,
	}
}

func RuleDelete(fnc DeleteFunc, paths ...string) Rule {
	return Rule{
		paths:      paths,
		maskFunc:   nil,
		deleteFunc: fnc,
	}
}

const Array string = "0..N"

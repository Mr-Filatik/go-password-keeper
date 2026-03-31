package dto

import (
	"errors"

	"github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask"
)

// Rules returns the rules for masking fields for the type.
//
// Implementation of the mask.IMaskable interface
// from the github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask package.
func (m UserInfo) Rules() []mask.Rule {
	return []mask.Rule{
		mask.RuleEdit(mask.Email(), "email"),
		mask.RuleEdit(mask.Password(), "password"),
		mask.RuleEdit(Custom(), "claims", mask.Array),
		mask.RuleEdit(Custom2(), "aaa", "cccs", mask.Array, "ccc"),
		mask.RuleEdit(Custom(), "aaa", "bbbs", mask.Array, "ccc"),
		mask.RuleDelete(mask.DAny(), "delete"),
	}
}

func Custom() mask.EditFunc {
	return mask.WrapStrToStrFn(func(s string) string {
		if s != "BBB" {
			return "*"
		}

		return s
	})
}

func Custom2() mask.EditFunc {
	return func(a any) (any, error) {
		return nil, errors.New("error")
	}
}

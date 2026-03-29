package dto

import (
	"github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask"
)

// Rules returns the rules for masking fields for the type.
//
// Implementation of the mask.IMaskable interface
// from the github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask package.
func (m UserInfo) Rules() []mask.Rule {
	return []mask.Rule{
		// UserInfo.Email
		mask.RuleEdit([]string{"email"}, mask.MaskedEmail()),
		// UserInfo.Password
		mask.RuleEdit([]string{"password"}, mask.MaskedPassword()),
		// UserInfo.Claims.0
		mask.RuleEditInSlice([]string{"claims"}, []string{}, Custom()),
		// UserInfo.AAA.BBB.0.CCC
		mask.RuleEditInSlice([]string{"aaa", "bbbs"}, []string{"ccc"}, Custom()),
	}
}

func Custom() mask.MFunc {
	return mask.WrapStrToStrFn(func(s string) string {
		if s != "BBB" {
			return "*"
		}

		return s
	})
}

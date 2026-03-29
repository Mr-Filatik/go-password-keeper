package dto

import "github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/sanitize"

type UserInfo struct {
	ID       string   `json:"id"       validate:"required"`
	Email    string   `json:"email"    validate:"required"`
	Password string   `json:"password" validate:"required"`
	Claims   []string `json:"claims"   validate:"required"`
	AAA      AAA      `json:"aaa"      validate:"required"`
}

type AAA struct {
	BBB []BBB `json:"bbbs" validate:"required"`
}

type BBB struct {
	CCC string `json:"ccc" validate:"required"`
}

func (m UserInfo) Sanitize() any {
	return UserInfo{
		ID:       m.ID,
		Email:    sanitize.StringMS(m.Email, sanitize.TypeMiddle, 4),
		Password: sanitize.StringM(m.Password, sanitize.TypeFull, 0),
	}
}

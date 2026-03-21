package dto

import "github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/sanitize"

type UserInfo struct {
	ID       string
	Email    string
	Password string
}

func (m UserInfo) Sanitize() any {
	return UserInfo{
		ID:       m.ID,
		Email:    sanitize.StringMS(m.Email, sanitize.TypeMiddle, 4),
		Password: sanitize.StringM(m.Password, sanitize.TypeFull, 0),
	}
}

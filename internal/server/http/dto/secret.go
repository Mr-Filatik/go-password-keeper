package dto

import (
	"fmt"

	"github.com/mr-filatik/go-password-keeper/internal/platform/types"
)

type Secret struct {
	ID   string     `json:"id"       validate:"required"`
	Name string     `json:"id"       validate:"required"`
	Type SecretType `json:"id"       validate:"required"`
}

// SecretType описывает тип секрета.
type SecretType string // @name SecretType

// type SecretTypeStrct struct {
// 	value string
// }

const (
	// Точные значения для типа SecretType в формате строки (упорядочены по алфавиту).
	nameSecretTypeJSON     = "json"
	nameSecretTypePassword = "password"
	nameSecretTypeString   = "string"

	// SecretTypeJSON описыает секрет имеющий формат JSON.
	SecretTypeJSON SecretType = nameSecretTypeJSON

	// SecretTypePassword описыает секрет имеющий формат скрытой строки.
	SecretTypePassword SecretType = nameSecretTypePassword

	// SecretTypeString описыает секрет имеющий формат строки.
	SecretTypeString SecretType = nameSecretTypeString
)

var ErrUnexpectedValueInEnumSecretType = fmt.Errorf("%w (SecretType)", types.ErrUnexpectedValueInEnum)

//nolint:gochecknoglobals
var mapOfSecretTypes = map[string]SecretType{
	nameSecretTypeJSON:     SecretTypeJSON,
	nameSecretTypePassword: SecretTypePassword,
	nameSecretTypeString:   SecretTypeString,
}

// SecretTypeFromString преобразует строку к типу перечисления.
func SecretTypeFromString(value string) (SecretType, error) {
	stValue, ok := mapOfSecretTypes[value]
	if ok {
		return stValue, nil
	}

	return SecretType(""), fmt.Errorf("%w: %q", ErrUnexpectedValueInEnumSecretType, value)
}

// analize memory go build -gcflags="-m" main.go

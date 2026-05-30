package types

import (
	"errors"
	"fmt"
)

var (
	// ErrUnexpectedValue описывает общую ошибку при получении неожиданного значения.
	//
	// В ошибках для конкретного типа указывать тип в следующем виде "%w in enum".
	// Пример: fmt.Errorf("%w in enum", ErrUnexpectedValue).
	ErrUnexpectedValue = errors.New("unexpected value")

	// ErrUnexpectedValueInEnum описывает общую ошибку при получении неожиданного значения в перечислении.
	//
	// В ошибках для конкретного перечисления указывать его название в следующем виде "%w (MyEnum)".
	// Пример: fmt.Errorf("%w (MyEnum)", types.ErrUnexpectedValueInEnum).
	ErrUnexpectedValueInEnum = fmt.Errorf("%w in enum", ErrUnexpectedValue)
)

var (
	// ErrParsing описывает общую ошибку при парсингe.
	ErrParsing = errors.New("failed to parse")
)

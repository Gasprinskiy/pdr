package shared

import (
	"errors"
)

var (
	ErrFileRead     = errors.New("an error ocured while file read")
	ErrLocalStorage = errors.New("app local storage error")
)

type LangCode int8

const (
	LangCodeRU LangCode = iota
	LangCodeEN
)

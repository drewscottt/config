package config

import (
	"errors"
)

var (
	ErrInvalidType                      = errors.New("invalid type")
	ErrValidation                       = errors.New("validaion error")
	ErrInvalidValueForType              = errors.New("invalid value for type")
	ErrCannotPutConfigTagOnPrivateField = errors.New("cannot put config tag on private field")
)

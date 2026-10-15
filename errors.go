package config

import (
	"errors"
)

var (
	ErrInvalidVariableType              = errors.New("invalid variable type")
	ErrValidation                       = errors.New("validaion error")
	ErrInvalidValueForVariableType      = errors.New("invalid value for variable type")
	ErrCannotPutConfigTagOnPrivateField = errors.New("cannot put config tag on private field")
	ErrConfigMustBePointerToStruct      = errors.New("config must be a pointer to a struct")
	ErrLoaderNotInitialized             = errors.New("loader not initialized")
)

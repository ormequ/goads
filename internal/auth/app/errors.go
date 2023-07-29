package app

import "errors"

var (
	ErrIncorrectCredentials = errors.New("incorrect credentials")
	ErrNotFound             = errors.New("not found")
	ErrInvalidContent       = errors.New("invalid content")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidToken         = errors.New("invalid token")
	ErrPasswordToShort      = errors.New("password to short")
)

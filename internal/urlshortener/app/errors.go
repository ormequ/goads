package app

import "errors"

var (
	ErrAlreadyExists   = errors.New("alias already exists")
	ErrNotFound        = errors.New("not found")
	ErrNoAds           = errors.New("ads not found")
	ErrAdNotExists     = errors.New("ad does not exist")
	ErrPermisionDenied = errors.New("permission denied")
)

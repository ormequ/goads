package app

import "errors"

var (
	ErrAdAlreadyAdded   = errors.New("ad has already been added")
	ErrAlreadyExists    = errors.New("alias or url already exists")
	ErrNotFound         = errors.New("not found")
	ErrNoAds            = errors.New("ads not found")
	ErrAdNotExists      = errors.New("ad does not exist")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidContent   = errors.New("invalid content")
)

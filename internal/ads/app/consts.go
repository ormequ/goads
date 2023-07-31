package app

import (
	"errors"
)

var (
	ErrPermissionDenied = errors.New("this user cannot edit the ad")
	ErrAdNotFound       = errors.New("ad not found")
	ErrAuthorNotFound   = errors.New("author not found")
	ErrInvalidContent   = errors.New("invalid ad's content")
	ErrInvalidFilter    = errors.New("invalid ad's filter")
)

const ServiceName = "Ads"

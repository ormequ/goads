package services

import (
	"errors"
	"goads/internal/ads/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetErrorStatus(err error) error {
	if err == nil {
		return nil
	}
	code := codes.Internal
	if errors.Is(err, app.ErrNotFound) {
		code = codes.NotFound
	}
	if errors.Is(err, app.ErrPermissionDenied) {
		code = codes.PermissionDenied
	}
	if errors.Is(err, app.ErrInvalidContent) || errors.Is(err, app.ErrInvalidFilter) {
		code = codes.InvalidArgument
	}
	return status.Error(code, err.Error())
}

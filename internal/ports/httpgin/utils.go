package httpgin

import (
	"errors"
	"goads/internal/app"
	"net/http"
)

func getErrorHTTPStatus(err error) int {
	if errors.Is(err, app.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, app.ErrPermissionDenied) {
		return http.StatusForbidden
	}
	if errors.Is(err, app.ErrInvalidContent) || errors.Is(err, app.ErrInvalidFilter) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

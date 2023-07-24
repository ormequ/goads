package responses

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goads/internal/app"
	"net/http"
)

func GetErrorHTTPStatus(err error) int {
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

func Error(err error) gin.H {
	if GetErrorHTTPStatus(err) == http.StatusInternalServerError {
		err = errors.New("internal server error")
	}
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

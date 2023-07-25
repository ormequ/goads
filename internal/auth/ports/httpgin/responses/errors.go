package responses

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/app"
	"net/http"
)

func GetErrorHTTPStatus(err error) int {
	if errors.Is(err, app.ErrIncorrectCredentials) || errors.Is(err, app.ErrInvalidToken) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, app.ErrEmailAlreadyExists) {
		return http.StatusConflict
	}
	if errors.Is(err, app.ErrInvalidContent) {
		return http.StatusBadRequest
	}
	if errors.Is(err, app.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func HiddenError(err error) gin.H {
	if GetErrorHTTPStatus(err) == http.StatusInternalServerError {
		//err = errors.New("internal server error")
	}
	return Error(err)
}

func Error(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

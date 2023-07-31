package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func GetHTTPStatus(err error) int {
	switch status.Code(err) {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadGateway
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.NotFound:
		return http.StatusNotFound
	}
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func HiddenResponse(err error) gin.H {
	if GetHTTPStatus(err) == http.StatusInternalServerError {
		err = errors.New("internal server error")
	}
	return Response(err)
}

func Response(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func ProceedResult(c *gin.Context, response gin.H, err error) {
	code := GetHTTPStatus(err)
	if code == http.StatusOK {
		c.JSON(code, response)
	} else {
		c.JSON(code, HiddenResponse(err))
	}
}

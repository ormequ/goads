package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/responses"
	"goads/internal/api/errors"
	"goads/internal/auth/proto"
	"net/http"
)

func Login(a proto.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req proto.AuthenticateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		token, err := a.Authenticate(c, &req)
		errors.ProceedResult(c, responses.TokenSuccess(token), err)
	}
}

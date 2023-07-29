package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/responses"
	"goads/internal/api/errors"
	"goads/internal/auth/proto"
	"net/http"
)

func Register(a proto.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req proto.RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		reg, err := a.Register(c, &req)
		errors.ProceedResult(c, responses.UserWithTokenSuccess(reg), err)
	}
}

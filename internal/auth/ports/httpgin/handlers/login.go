package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/ports/httpgin/responses"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authenticator interface {
	Authenticate(ctx context.Context, email string, password string) (string, error)
}

func Login(a authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		token, err := a.Authenticate(c, req.Email, req.Password)
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.TokenSuccess(token))
	}
}

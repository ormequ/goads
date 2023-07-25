package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/ports/httpgin/responses"
	"goads/internal/auth/users"
	"net/http"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerer interface {
	authenticator
	Register(ctx context.Context, email string, name string, password string) (users.User, error)
}

func Register(a registerer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := a.Register(c, req.Email, req.Name, req.Password)
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		token, err := a.Authenticate(c, req.Email, req.Password)
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.UserWithTokenSuccess(user, token))
	}
}

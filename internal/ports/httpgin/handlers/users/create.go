package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/entities/users"
	"goads/internal/ports/httpgin/responses"
	"net/http"
)

type creator interface {
	Create(ctx context.Context, email string, name string) (users.User, error)
}

type createRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

func Create(app creator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := app.Create(c, req.Email, req.Name)
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

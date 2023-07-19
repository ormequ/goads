package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/entities/users"
	"goads/internal/ports/httpgin/responses"
	"net/http"
	"strconv"
)

type getterByID interface {
	GetByID(ctx context.Context, id int64) (users.User, error)
}

type getterByEmail interface {
	GetByEmail(ctx context.Context, email string) (users.User, error)
}

type getByEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func GetByID(app getterByID) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := app.GetByID(c, int64(id))
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

func GetByEmail(app getterByEmail) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getByEmailRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := app.GetByEmail(c, req.Email)
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/entities/users"
	"goads/internal/ports/httpgin/responses"
	"net/http"
	"strconv"
)

type emailChanger interface {
	getterByID
	ChangeEmail(ctx context.Context, id int64, email string) error
}

type nameChanger interface {
	getterByID
	ChangeName(ctx context.Context, id int64, name string) error
}

type changeUserEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

type changeUserNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func ChangeEmail(app emailChanger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserEmailRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = app.ChangeEmail(c, int64(id), req.Email)
		var user users.User
		if err == nil {
			user, err = app.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

func ChangeName(app nameChanger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserNameRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = app.ChangeName(c, int64(id), req.Name)
		var user users.User
		if err == nil {
			user, err = app.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/ports/httpgin/responses"
	"goads/internal/auth/ports/httpgin/utils"
	"goads/internal/auth/users"
	"net/http"
	"strconv"
)

type getterByID interface {
	GetByID(ctx context.Context, id int64) (users.User, error)
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
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

type validator interface {
	Validate(ctx context.Context, token string) (users.User, error)
}

func GetFromToken(app validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := app.Validate(c, utils.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

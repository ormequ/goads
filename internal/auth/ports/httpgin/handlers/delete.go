package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/ports/httpgin/responses"
	"goads/internal/auth/ports/httpgin/utils"
	"net/http"
)

type deleter interface {
	validator
	Delete(ctx context.Context, id int64) error
}

func Delete(a deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := a.Validate(c, utils.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.Error(err))
			return
		}
		err = a.Delete(c, user.ID)
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.EmptySuccess())
	}
}

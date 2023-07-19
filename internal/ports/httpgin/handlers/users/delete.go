package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/ports/httpgin/responses"
	"net/http"
	"strconv"
)

type deleter interface {
	Delete(ctx context.Context, id int64) error
}

func Delete(a deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = a.Delete(c, int64(id))
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.EmptySuccess())
	}
}

package ads

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/ports/httpgin/responses"
	"net/http"
	"strconv"
)

type deleter interface {
	Delete(ctx context.Context, id int64, userID int64) error
}

type deleteAdRequest struct {
	UserID int64 `json:"user_id"`
}

func Delete(app deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req deleteAdRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = app.Delete(c, int64(id), req.UserID)
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.EmptySuccess())
	}
}

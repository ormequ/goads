package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/ads/ads"
	responses2 "goads/internal/ads/ports/httpgin/responses"
	"net/http"
	"strconv"
)

type statusChanger interface {
	getterByID
	ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error
}

type changeStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

func ChangeStatus(app statusChanger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeStatusRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses2.Error(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses2.Error(err))
			return
		}
		err = app.ChangeStatus(c, int64(id), req.UserID, req.Published)
		var ad ads.Ad
		if err == nil {
			ad, err = app.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(responses2.GetErrorHTTPStatus(err), responses2.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses2.AdSuccess(ad))
	}
}

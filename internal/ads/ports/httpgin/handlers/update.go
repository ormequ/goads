package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/ads/ads"
	responses2 "goads/internal/ads/ports/httpgin/responses"
	"net/http"
	"strconv"
)

type updater interface {
	getterByID
	Update(ctx context.Context, id int64, userID int64, title string, text string) error
}

type updateRequest struct {
	Title  string `json:"title" binding:"required"`
	Text   string `json:"text" binding:"required"`
	UserID int64  `json:"user_id"`
}

func Update(app updater) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses2.Error(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses2.Error(err))
			return
		}
		err = app.Update(c, int64(id), req.UserID, req.Title, req.Text)
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

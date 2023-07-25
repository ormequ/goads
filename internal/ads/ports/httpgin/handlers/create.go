package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"goads/internal/ads/ads"
	responses2 "goads/internal/ads/ports/httpgin/responses"
	"net/http"
)

type creator interface {
	Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error)
}

type createRequest struct {
	Title  string `json:"title" binding:"required"`
	Text   string `json:"text" binding:"required"`
	UserID int64  `json:"user_id"`
}

func Create(app creator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses2.Error(err))
			return
		}
		ad, err := app.Create(c, req.Title, req.Text, req.UserID)
		if err != nil {
			c.JSON(responses2.GetErrorHTTPStatus(err), responses2.Error(err))
			return
		}
		fmt.Println(ad)
		c.JSON(http.StatusOK, responses2.AdSuccess(ad))
	}
}

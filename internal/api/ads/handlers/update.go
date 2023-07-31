package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/ads/proto"
	"goads/internal/api/ads/responses"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"net/http"
	"strconv"
)

type updateRequest struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
}

func Update(client proto.AdServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		userID, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		ad, err := client.Update(c, &proto.UpdateAdRequest{
			AdId:     int64(id),
			Title:    req.Title,
			Text:     req.Text,
			AuthorId: userID,
		})
		errors.ProceedResult(c, responses.AdSuccess(ad), err)
	}
}

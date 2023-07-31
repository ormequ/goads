package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/api/urlshortener/responses"
	"goads/internal/urlshortener/proto"
	"net/http"
)

type createRequest struct {
	URL   string  `json:"url" binding:"required"`
	Alias string  `json:"alias"`
	Ads   []int64 `json:"ads"`
}

func Create(shortener proto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createRequest
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		userID, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		link, err := shortener.Create(c, &proto.CreateRequest{
			Url:      req.URL,
			Alias:    req.Alias,
			AuthorId: userID,
			Ads:      req.Ads,
		})
		errors.ProceedResult(c, responses.LinkSuccess(link), err)
	}
}

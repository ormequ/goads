package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/ads/proto"
	"goads/internal/api/ads/responses"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"net/http"
	"time"
)

func GetFiltered(client proto.AdServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := &proto.FilterAdsRequest{}
		if _, ok := c.GetQuery("all"); ok {
			filter.All = true
		}
		if date, ok := c.GetQuery("date"); ok {
			datetime, err := time.Parse(time.RFC3339Nano, date)
			if err != nil {
				c.JSON(http.StatusBadRequest, errors.Response(err))
				return
			}
			filter.Date = datetime.UnixMilli()
		}
		if q, ok := c.GetQuery("search"); ok {
			filter.Title = q
		}
		var err error
		filter.AuthorId, err = utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		adsList, err := client.Filter(c, filter)
		errors.ProceedResult(c, responses.AdsSuccess(adsList), err)
	}
}

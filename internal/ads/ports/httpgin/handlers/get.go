package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/ads/ads"
	"goads/internal/ads/app"
	responses2 "goads/internal/ads/ports/httpgin/responses"
	"net/http"
	"strconv"
	"time"
)

type getterByID interface {
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
}

type getterFiltered interface {
	GetFiltered(ctx context.Context, filter app.Filter) ([]ads.Ad, error)
}

func GetFiltered(a getterFiltered) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := app.Filter{
			AuthorID: -1,
		}
		for _, key := range [3]string{"all", "date", "author"} {
			val, ok := c.GetQuery(key)
			if !ok {
				continue
			}
			switch key {
			case "all":
				filter.All = true
			case "author":
				id, err := strconv.Atoi(val)
				if err != nil {
					c.JSON(http.StatusBadRequest, responses2.Error(err))
					return
				}
				filter.AuthorID = int64(id)
			case "date":
				date, err := time.Parse(time.RFC3339Nano, val)
				if err != nil {
					c.JSON(http.StatusBadRequest, responses2.Error(err))
					return
				}
				filter.Date = date
			}
		}
		adsList, err := a.GetFiltered(c, filter)
		if err != nil {
			c.JSON(responses2.GetErrorHTTPStatus(err), responses2.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses2.AdsListSuccess(adsList))
	}
}

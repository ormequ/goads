package ads

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/app/ad"
	"goads/internal/entities/ads"
	"goads/internal/ports/httpgin/responses"
	"net/http"
	"strconv"
	"time"
)

type getterByID interface {
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
}

type getterFiltered interface {
	GetFiltered(ctx context.Context, filter ad.Filter) ([]ads.Ad, error)
}

func GetFiltered(app getterFiltered) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := ad.Filter{
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
					c.JSON(http.StatusBadRequest, responses.Error(err))
					return
				}
				filter.AuthorID = int64(id)
			case "date":
				date, err := time.Parse(time.RFC3339Nano, val)
				if err != nil {
					c.JSON(http.StatusBadRequest, responses.Error(err))
					return
				}
				filter.Date = date
			}
		}
		adsList, err := app.GetFiltered(c, filter)
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.AdsListSuccess(adsList))
	}
}

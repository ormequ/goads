package ads

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/entities/ads"
	"goads/internal/ports/httpgin/responses"
	"net/http"
)

type searcher interface {
	Search(ctx context.Context, title string) ([]ads.Ad, error)
}

func Search(app searcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		adsList, err := app.Search(c, c.Query("q"))
		if err != nil {
			c.JSON(responses.GetErrorHTTPStatus(err), responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.AdsListSuccess(adsList))
	}
}

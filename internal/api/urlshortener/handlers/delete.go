package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/api/urlshortener/responses"
	"goads/internal/urlshortener/proto"
	"net/http"
	"strconv"
)

func Delete(shortener proto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		_, err = shortener.Delete(c, &proto.DeleteRequest{
			Id:       int64(id),
			AuthorId: userID,
		})
		errors.ProceedResult(c, responses.EmptySuccess(), err)
	}
}

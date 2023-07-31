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

func Delete(client proto.AdServiceClient) gin.HandlerFunc {
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
		_, err = client.Delete(c, &proto.DeleteAdRequest{
			AdId:     int64(id),
			AuthorId: userID,
		})
		errors.ProceedResult(c, responses.EmptySuccess(), err)
	}
}

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

type changeStatusRequest struct {
	Published bool `json:"published" binding:"required"`
}

func ChangeStatus(client proto.AdServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeStatusRequest
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
		ad, err := client.ChangeStatus(c, &proto.ChangeAdStatusRequest{
			AdId:      int64(id),
			AuthorId:  userID,
			Published: req.Published,
		})
		errors.ProceedResult(c, responses.AdSuccess(ad), err)
	}
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/responses"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/auth/proto"
)

func GetFromToken(app proto.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		user, err := app.GetByID(c, &proto.GetUserByIDRequest{Id: id})
		errors.ProceedResult(c, responses.UserSuccess(user), err)
	}
}

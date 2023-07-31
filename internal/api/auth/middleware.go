package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/auth/proto"
	"net/http"
)

func Middleware(client proto.AuthServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := client.Validate(c, &proto.ValidateRequest{Token: utils.ExtractToken(c)})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.Response(fmt.Errorf("invalid token")))
			return
		}
		c.Set("userID", user.Id)
		c.Next()
	}
}

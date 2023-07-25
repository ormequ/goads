package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/app"
	"goads/internal/auth/ports/httpgin/responses"
	"goads/internal/auth/ports/httpgin/utils"
	"net/http"
)

func AuthorizeMW(v app.Validator) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, err := v.Validate(c, utils.ExtractToken(c))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.Error(errors.New("invalid token")))
			return
		}
		c.Next()
	}
}

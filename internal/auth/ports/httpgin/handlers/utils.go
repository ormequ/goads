package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/auth/ports/httpgin/responses"
)

func displayHiddenErr(c *gin.Context, err error) {
	c.JSON(responses.GetErrorHTTPStatus(err), responses.HiddenError(err))
}

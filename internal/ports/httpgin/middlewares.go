package httpgin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// LoggerMW logs information about all requests into out
func LoggerMW(c *gin.Context) {
	start := time.Now()

	c.Next()

	url := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	if query != "" {
		url += "?" + query
	}
	errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
	log.Printf("[LOG] %d | %13v | %15s | %-7s %#v\n%s",
		c.Writer.Status(),
		time.Since(start),
		c.ClientIP(),
		c.Request.Method,
		url,
		errMsg,
	)
}

// RecoveryMW recovers panics and writes information about them into out
func RecoveryMW(c *gin.Context) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		log.Printf("[Recovery] Panic recovered:\n%s\n", string(debug.Stack()))
		c.AbortWithStatus(http.StatusInternalServerError)
	}()
	c.Next()
}

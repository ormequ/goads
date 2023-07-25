package httpgin

import (
	"github.com/gin-gonic/gin"
	"goads/internal/ads/app"
	handlers2 "goads/internal/ads/ports/httpgin/handlers"
)

func SetRoutes(r *gin.RouterGroup, a app.App) {
	r.GET("/ads", handlers2.GetFiltered(a))                // Get all ads with filters applying
	r.GET("/ads/search", handlers2.Search(a))              // Search in ads
	r.POST("/ads", handlers2.Create(a))                    // Register an ad
	r.PUT("/ads/:ad_id/status", handlers2.ChangeStatus(a)) // Change ad's status [(un-)publishing]
	r.PUT("/ads/:ad_id", handlers2.Update(a))              // Change ad's content [title and text]
	r.DELETE("/ads/:ad_id", handlers2.Delete(a))           // Delete ad

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

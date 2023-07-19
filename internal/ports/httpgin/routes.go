package httpgin

import (
	"github.com/gin-gonic/gin"
	"goads/internal/app/ad"
	"goads/internal/app/user"
	"goads/internal/ports/httpgin/handlers/ads"
	"goads/internal/ports/httpgin/handlers/users"
)

func SetRoutes(r *gin.RouterGroup, a ad.App, u user.App) {
	r.GET("/ads", ads.GetFiltered(a))                // Get all ads with filters applying
	r.GET("/ads/search", ads.Search(a))              // Search in ads
	r.POST("/ads", ads.Create(a))                    // Create an ad
	r.PUT("/ads/:ad_id/status", ads.ChangeStatus(a)) // Change ad's status [(un-)publishing]
	r.PUT("/ads/:ad_id", ads.Update(a))              // Change ad's content [title and text]
	r.DELETE("/ads/:ad_id", ads.Delete(a))           // Delete ad

	r.GET("/users/:user_id", users.GetByID(u))           // Get user by id
	r.POST("/users", users.GetByEmail(u))                // Get user by email
	r.POST("/users/create", users.Create(u))             // Create a user
	r.PUT("/users/:user_id/name", users.ChangeName(u))   // Change user's name
	r.PUT("/users/:user_id/email", users.ChangeEmail(u)) // Change user's email
	r.DELETE("/users/:user_id", users.Delete(u))         // Delete user

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

package httpgin

import (
	"github.com/gin-gonic/gin"
	"goads/internal/app"
)

func SetRoutes(r *gin.RouterGroup, a app.Ads, u app.Users) {
	r.GET("/ads", getAdsFiltered(a))               // Get all ads with filters applying
	r.GET("/ads/search", searchAds(a))             // Search in ads
	r.POST("/ads", createAd(a))                    // Create an ad
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Change ad's status [(un-)publishing]
	r.PUT("/ads/:ad_id", updateAd(a))              // Change ad's content [title and text]
	r.DELETE("/ads/:ad_id", deleteAd(a))           // Delete ad

	r.GET("/users/:user_id", getUser(u))               // Get user by id
	r.POST("/users", createUser(u))                    // Create a user
	r.PUT("/users/:user_id/name", changeUserName(u))   // Change user's name
	r.PUT("/users/:user_id/email", changeUserEmail(u)) // Change user's email
	r.DELETE("/users/:user_id", deleteUser(u))         // Delete user

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

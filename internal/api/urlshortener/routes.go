package urlshortener

import (
	"github.com/gin-gonic/gin"
	adProto "goads/internal/ads/proto"
	"goads/internal/api/auth"
	"goads/internal/api/urlshortener/handlers"
	authProto "goads/internal/auth/proto"
	shProto "goads/internal/urlshortener/proto"
)

func SetRoutes(r gin.IRouter, authSvc authProto.AuthServiceClient, shortener shProto.ShortenerServiceClient, ads adProto.AdServiceClient) {
	r.GET("link/:alias", handlers.GetByAlias(shortener, ads))

	links := r.Group("/links")
	links.Use(auth.Middleware(authSvc))
	links.GET("/", handlers.GetByAuthor(shortener))
	links.POST("/", handlers.Create(shortener))
	links.GET("/:link_id", handlers.GetByID(shortener))
	links.PUT("/:link_id", handlers.UpdateAlias(shortener))
	links.DELETE("/:link_id", handlers.Delete(shortener))
	links.PUT("/:link_id/ads", handlers.UpdateAdData(shortener.AddAd))
	links.DELETE("/:link_id/ads", handlers.UpdateAdData(shortener.DeleteAd))
}

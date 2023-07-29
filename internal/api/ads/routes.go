package ads

import (
	"github.com/gin-gonic/gin"
	"goads/internal/ads/proto"
	"goads/internal/api/ads/handlers"
	"goads/internal/api/auth"
	authProto "goads/internal/auth/proto"
)

func SetRoutes(r gin.IRouter, authSvc authProto.AuthServiceClient, client proto.AdServiceClient) {
	g := r.Group("/ads")
	g.Use(auth.Middleware(authSvc))
	g.GET("/", handlers.GetFiltered(client))
	g.POST("/", handlers.Create(client))
	g.PUT("/:ad_id/status", handlers.ChangeStatus(client))
	g.PUT("/:ad_id", handlers.Update(client))
	g.DELETE("/:ad_id", handlers.Delete(client))
}

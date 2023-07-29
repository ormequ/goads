package auth

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/handlers"
	"goads/internal/auth/proto"
)

func SetRoutes(r gin.IRouter, client proto.AuthServiceClient) {
	r.POST("/register", handlers.Register(client))
	r.POST("/login", handlers.Login(client))

	auth := r.Group("/user")
	auth.Use(Middleware(client))
	auth.GET("/", handlers.GetFromToken(client))
	auth.PUT("/name", handlers.ChangeName(client))
	auth.PUT("/email", handlers.ChangeEmail(client))
	auth.PUT("/password", handlers.ChangePassword(client))
	auth.DELETE("/", handlers.Delete(client))
}

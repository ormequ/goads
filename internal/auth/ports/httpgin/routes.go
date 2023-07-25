package httpgin

import (
	"github.com/gin-gonic/gin"
	"goads/internal/auth/app"
	"goads/internal/auth/ports/httpgin/handlers"
)

func SetRoutes(r *gin.RouterGroup, a app.App) {
	r.GET("/users/:user_id", handlers.GetByID(a))
	r.POST("/register", handlers.Register(a))
	r.POST("/login", handlers.Login(a))

	auth := r.Group("/")
	auth.Use(AuthorizeMW(a.Validator))
	auth.GET("/user", handlers.GetFromToken(a))
	auth.PUT("/user/name", handlers.ChangeName(a))
	auth.PUT("/user/email", handlers.ChangeEmail(a))
	auth.PUT("/user/password", handlers.ChangePassword(a))
	auth.DELETE("/user", handlers.Delete(a))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

package server

import (
	"context"
	"errors"
	"fmt"
	adProto "goads/internal/ads/proto"
	"goads/internal/api/ads"
	"goads/internal/api/auth"
	"goads/internal/api/urlshortener"
	authProto "goads/internal/auth/proto"
	shProto "goads/internal/urlshortener/proto"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	http.Server
}

func New(addr string, authSvc authProto.AuthServiceClient, shSvc shProto.ShortenerServiceClient, adsSvc adProto.AdServiceClient) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	s := Server{http.Server{
		Addr:    addr,
		Handler: r,
	}}
	api := r.Group("/api")
	api.Any("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	auth.SetRoutes(api, authSvc)
	ads.SetRoutes(api, authSvc, adsSvc)
	urlshortener.SetRoutes(api, authSvc, shSvc, adsSvc)
	return &s
}

func (s *Server) Listen(ctx context.Context) error {
	errCh := make(chan error)
	defer func() {
		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(shCtx); err != nil {
			log.Printf("can't close http server listening on %s: %s", s.Addr, err.Error())
		}
		close(errCh)
	}()

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return fmt.Errorf("http server can't listen and serve requests: %w", err)
	}
}

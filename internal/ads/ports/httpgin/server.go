package httpgin

import (
	"context"
	"errors"
	"fmt"
	"goads/internal/ads/app"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	http.Server
}

func NewServer(addr string, a app.App) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerMW)
	r.Use(RecoveryMW)
	s := Server{http.Server{
		Addr:    addr,
		Handler: r,
	}}
	api := r.Group("/api/v1")
	SetRoutes(api, a)
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

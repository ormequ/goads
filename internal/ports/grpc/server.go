package grpc

import (
	"context"
	"fmt"
	"goads/internal/app/ad"
	"goads/internal/app/user"
	"goads/internal/ports/grpc/services"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	port   string
	server *grpc.Server
}

func (s *Server) Listen(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	defer func() {
		s.server.GracefulStop()
		_ = lis.Close()
		close(errCh)
	}()
	go func() {
		if err := s.server.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return fmt.Errorf("grpc server error: %w", err)
	}
}

func NewServer(port string, a ad.App, u user.App) Server {
	s := grpc.NewServer(GetUnaryInterceptors())
	services.RegisterUserServiceServer(s, services.NewUsers(u))
	services.RegisterAdServiceServer(s, services.NewAds(a))
	return Server{port, s}
}

func GetUnaryInterceptors() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(LogUnaryInterceptor, RecoveryUnaryInterceptor())
}

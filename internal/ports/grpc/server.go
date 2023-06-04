package grpc

import (
	"context"
	"fmt"
	"goads/internal/app"
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

func NewServer(port string, a app.Ads, u app.Users) Server {
	s := grpc.NewServer(GetUnaryInterceptors())
	RegisterAdServiceServer(s, NewService(a, u))
	return Server{port, s}
}

func GetUnaryInterceptors() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(LogUnaryInterceptor, RecoveryUnaryInterceptor())
}

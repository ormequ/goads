package grpc

import (
	"context"
	"fmt"
	"goads/internal/ads/app"
	services2 "goads/internal/ads/ports/grpc/services"
	"goads/internal/ads/proto"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	addr   string
	server *grpc.Server
}

func (s Server) Listen(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
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

func NewServer(addr string, a app.App) Server {
	s := grpc.NewServer(GetUnaryInterceptors())
	proto.RegisterAdServiceServer(s, services2.NewAds(a))
	return Server{addr, s}
}

func GetUnaryInterceptors() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(LogUnaryInterceptor, RecoveryUnaryInterceptor())
}

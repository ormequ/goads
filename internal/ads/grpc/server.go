package grpc

import (
	"context"
	"fmt"
	"goads/internal/ads/app"
	"goads/internal/ads/proto"
	grpcInt "goads/internal/pkg/grpc"
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
	s := grpc.NewServer(grpcInt.GetUnaryInterceptors())
	proto.RegisterAdServiceServer(s, NewService(a))
	return Server{addr, s}
}

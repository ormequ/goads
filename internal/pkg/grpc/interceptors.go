package grpc

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func LogUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf(
		"[LOG] gRPC | Unary | %s\n",
		info.FullMethod,
	)
	return handler(ctx, req)
}

func recoveryHandler(p any) error {
	return status.Errorf(codes.Internal, "%v [RECOVERY] %v\n", time.Now().Format("2006/01/02 15:04:05"), p)
}

func RecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(recoveryHandler))
}

func GetUnaryInterceptors() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(LogUnaryInterceptor, RecoveryUnaryInterceptor())
}

package main

import (
	"context"
	"fmt"
	adProto "goads/internal/ads/proto"
	"goads/internal/api/server"
	authProto "goads/internal/auth/proto"
	"goads/internal/pkg/config"
	"goads/internal/pkg/shutdown"
	shProto "goads/internal/urlshortener/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

type Config struct {
	Env              string `env:"ENV" env-default:"local"`
	HTTPAddress      string `env:"HTTP_ADDRESS" env-default:":80"`
	URLShortenerPath string `env:"URL_SHORTENER_PATH" env-required:"true"`
	AuthPath         string `env:"AUTH_PATH" env-required:"true"`
	AdsPath          string `env:"ADS_PATH" env-required:"true"`
}

func connect(ctx context.Context, name string, path string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, path, grpc.WithTransportCredentials(insecure.NewCredentials()))
	for i := 0; i < 10 && err != nil; i++ {
		conn, err = grpc.DialContext(ctx, path, grpc.WithTransportCredentials(insecure.NewCredentials()))
		fmt.Printf("Reconnect to %s #%d", name, i+1)
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		log.Fatalf("Cannot start connection with %s: %v", name, err)
	}
	return conn
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))

	eg, ctx := errgroup.WithContext(context.Background())

	shConn := connect(ctx, "URL Shortener", cfg.URLShortenerPath)
	authConn := connect(ctx, "Auth", cfg.AuthPath)
	adsConn := connect(ctx, "Ads", cfg.AdsPath)

	shSvc := shProto.NewShortenerServiceClient(shConn)
	authSvc := authProto.NewAuthServiceClient(authConn)
	adsSvc := adProto.NewAdServiceClient(adsConn)

	srv := server.New(cfg.HTTPAddress, authSvc, shSvc, adsSvc)
	shutdown.Gracefully(eg, ctx, srv)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown server:", err)
	}
	log.Println("Server has been shutdown successfully")
}

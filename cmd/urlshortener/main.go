package main

import (
	"context"
	"fmt"
	adProto "goads/internal/ads/proto"
	"goads/internal/pkg/config"
	"goads/internal/pkg/shutdown"
	"goads/internal/urlshortener/adapters/ads"
	"goads/internal/urlshortener/adapters/pgrepo"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/generator"
	grpcPort "goads/internal/urlshortener/grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Env          string `env:"ENV" env-default:"local"`
	GRPCAddress  string `env:"GRPC_ADDRESS" env-default:":8081"`
	PostgresConn string `env:"POSTGRES_CONN" env-required:"true"`
	AdsPath      string `env:"ADS_PATH" env-required:"true"`
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))

	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(ctx, cfg.PostgresConn)
	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(time.Second * 3)
		fmt.Printf("Reconnect to PostgreSQL #%d", i+1)
		conn, err = pgx.Connect(ctx, cfg.PostgresConn)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close(ctx) }()
	adsConn, err := grpc.DialContext(ctx, cfg.AdsPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(time.Second * 3)
		fmt.Printf("Reconnect to Ads #%d", i+1)
		adsConn, err = grpc.DialContext(ctx, cfg.AdsPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if err != nil {
		log.Fatalf("Cannot start connection with Ads: %v", err)
	}
	adsSvc := adProto.NewAdServiceClient(adsConn)

	repo := pgrepo.New(conn, conn)
	a := app.New(repo, generator.New(repo), ads.New(adsSvc))

	grpcServer := grpcPort.NewServer(cfg.GRPCAddress, a)

	shutdown.Gracefully(eg, ctx, grpcServer)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown server:", err)
	}
	log.Println("Server has been shutdown successfully")
}

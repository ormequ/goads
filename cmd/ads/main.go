package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"goads/internal/ads/adapters/pgrepo"
	"goads/internal/ads/app"
	"goads/internal/ads/grpc"
	"goads/internal/pkg/config"
	"goads/internal/pkg/shutdown"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"
)

type Config struct {
	Env          string `env:"ENV" env-default:"local"`
	GRPCAddress  string `env:"GRPC_ADDRESS" env-default:":8888"`
	PostgresConn string `env:"POSTGRES_CONN" env-required:"true"`
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))
	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(ctx, cfg.PostgresConn)
	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(time.Second * 3)
		fmt.Printf("Reconnect to PostgreSQL #%d\n", i+1)
		conn, err = pgx.Connect(ctx, cfg.PostgresConn)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close(ctx) }()

	grpcServer := grpc.NewServer(cfg.GRPCAddress, app.New(pgrepo.New(conn)))

	shutdown.Gracefully(eg, ctx, grpcServer)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown servers:", err)
	}
	log.Println("Servers have been shutdown successfully")
}

package main

import (
	"context"
	"goads/internal/pkg/config"
	"goads/internal/pkg/shutdown"
	"goads/internal/urlshortener/adapters/pgrepo"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/generator"
	"goads/internal/urlshortener/grpc"
	"golang.org/x/sync/errgroup"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Env          string `env:"ENV" env-default:"local"`
	GRPCAddress  string `env:"GRPC_ADDRESS" env-default:":8081"`
	PostgresConn string `env:"POSTGRES_CONN" env-required:"true"`
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))

	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(ctx, cfg.PostgresConn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close(ctx) }()

	repo := pgrepo.New(conn, conn)
	a := app.New(repo, generator.New(repo))

	grpcServer := grpc.NewServer(cfg.GRPCAddress, a)

	shutdown.Gracefully(eg, ctx, grpcServer)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown server:", err)
	}
	log.Println("Server has been shutdown successfully")
}

package main

import (
	"context"
	"fmt"
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
	Env         string `env:"ENV" env-default:"local"`
	GRPCAddress string `env:"GRPC_ADDRESS" env-default:":8081"`
	DBHost      string `env:"POSTGRES_HOST" env-required:"true"`
	DBPort      uint16 `env:"POSTGRES_PORT" env-required:"true"`
	DBUser      string `env:"POSTGRES_USER" env-required:"true"`
	DBPassword  string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName      string `env:"POSTGRES_DB" env-required:"true"`
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))

	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
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

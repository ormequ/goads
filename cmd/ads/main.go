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
)

type Config struct {
	Env              string `env:"ENV" env-default:"local"`
	GRPCAddress      string `env:"GRPC_ADDRESS" env-default:":18081"`
	PostgresHost     string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort     uint16 `env:"POSTGRES_PORT" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))
	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(
		ctx,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s", cfg.PostgresUser,
			cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB,
		),
	)
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

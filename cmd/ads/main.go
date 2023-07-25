package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"goads/internal/ads/adapters/pgrepo"
	"goads/internal/ads/app"
	"goads/internal/ads/config"
	"goads/internal/ads/ports/grpc"
	"goads/internal/ads/ports/httpgin"
	"goads/internal/pkg/shutdown"
	"golang.org/x/sync/errgroup"
	"log"
)

func main() {
	cfg := config.MustLoad()
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

	adsApp := app.New(pgrepo.New(conn))
	httpServer := httpgin.NewServer(cfg.HTTPAddress, adsApp)
	grpcServer := grpc.NewServer(cfg.GRPCAddress, adsApp)

	shutdown.SetupGraceful(eg, ctx, httpServer, grpcServer)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown servers:", err)
	}
	log.Println("Servers have been shutdown successfully")
}

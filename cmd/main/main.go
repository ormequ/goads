package main

import (
	"context"
	"fmt"
	"goads/internal/adapters/maprepo"
	"goads/internal/app/ad"
	"goads/internal/app/user"
	"goads/internal/config"
	"goads/internal/ports/grpc"
	"goads/internal/ports/httpgin"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
)

func main() {
	cfg := config.MustLoad()
	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close(ctx) }()

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	usersProv := user.New(maprepo.NewUsers())
	adsProv := ad.New(maprepo.NewAds(), usersProv)
	httpServer := httpgin.NewServer(cfg.HTTPAddress, adsProv, usersProv)
	grpcServer := grpc.NewServer(cfg.GRPCAddress, adsProv, usersProv)

	eg.Go(func() error {
		return httpServer.Listen(ctx)
	})
	eg.Go(func() error {
		return grpcServer.Listen(ctx)
	})

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown servers:", err)
	}
	log.Println("Servers have been shutdown successfully")
}

package main

import (
	"context"
	"fmt"
	"goads/internal/adapters/maprepo"
	"goads/internal/app/providers"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"goads/internal/ports/grpc"
	"goads/internal/ports/httpgin"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	eg, ctx := errgroup.WithContext(context.Background())

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

	httpPort := ":18080"
	grpcPort := ":18081"

	usersProv := providers.NewUsers(maprepo.New[users.User]())
	adsProv := providers.NewAds(maprepo.New[ads.Ad](), usersProv)
	httpServer := httpgin.NewServer(httpPort, adsProv, usersProv)
	grpcServer := grpc.NewServer(grpcPort, adsProv, usersProv)

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

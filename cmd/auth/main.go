package main

import (
	"context"
	"fmt"
	"goads/internal/auth/adapters/bcrypt"
	"goads/internal/auth/adapters/jwt"
	"goads/internal/auth/adapters/pgrepo"
	"goads/internal/auth/app"
	"goads/internal/auth/config"
	"goads/internal/auth/ports/grpc"
	"goads/internal/auth/ports/httpgin"
	"goads/internal/pkg/shutdown"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func mustReadFile(file string) []byte {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func main() {
	cfg := config.MustLoad()
	tokenizer, err := jwt.NewTokenizer(time.Duration(cfg.Expires)*time.Hour, mustReadFile(cfg.PrivateKey))
	if err != nil {
		log.Fatal(err)
	}
	validator, err := jwt.NewValidator(mustReadFile(cfg.PublicKey))
	if err != nil {
		log.Fatal(err)
	}

	eg, ctx := errgroup.WithContext(context.Background())

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close(ctx) }()

	a := app.New(pgrepo.New(conn), tokenizer, bcrypt.New(cfg.PasswordCost), validator)

	httpServer := httpgin.NewServer(cfg.HTTPAddress, a)
	grpcServer := grpc.NewServer(cfg.GRPCAddress, a)

	shutdown.SetupGraceful(eg, ctx, httpServer, grpcServer)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown server:", err)
	}
	log.Println("Server has been shutdown successfully")
}

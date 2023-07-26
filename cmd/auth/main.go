package main

import (
	"context"
	"fmt"
	"goads/internal/auth/adapters/bcrypt"
	"goads/internal/auth/adapters/jwt"
	"goads/internal/auth/adapters/pgrepo"
	"goads/internal/auth/app"
	"goads/internal/auth/ports/grpc"
	"goads/internal/auth/ports/httpgin"
	"goads/internal/pkg/config"
	"goads/internal/pkg/shutdown"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Env          string `env:"ENV" env-default:"local"`
	PrivateKey   string `env:"AUTH_PRIVATE_KEY" env-required:"true"`
	PublicKey    string `env:"AUTH_PUBLIC_KEY" env-required:"true"`
	Expires      int    `env:"AUTH_EXPIRES_HOURS" env-default:"24"`
	PasswordCost int    `env:"PASSWORD_COST" env-default:"10"`
	HTTPAddress  string `env:"HTTP_ADDRESS" env-default:":8080"`
	GRPCAddress  string `env:"GRPC_ADDRESS" env-default:":8081"`
	DBHost       string `env:"POSTGRES_HOST" env-required:"true"`
	DBPort       uint16 `env:"POSTGRES_PORT" env-required:"true"`
	DBUser       string `env:"POSTGRES_USER" env-required:"true"`
	DBPassword   string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName       string `env:"POSTGRES_DB" env-required:"true"`
}

func mustReadFile(file string) []byte {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func main() {
	cfg := config.MustLoadENV[Config](os.Getenv("CONFIG_PATH"))
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

	shutdown.Gracefully(eg, ctx, httpServer, grpcServer)

	if err := eg.Wait(); err != nil {
		log.Println("Graceful shutdown server:", err)
	}
	log.Println("Server has been shutdown successfully")
}

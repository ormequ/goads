package bcrypt

import (
	"context"
	"errors"
	"goads/internal/auth/app"

	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	cost int
}

func (b BCrypt) Generate(ctx context.Context, password string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(bytes), err
}

func (b BCrypt) Compare(ctx context.Context, hash string, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return app.ErrIncorrectCredentials
	}
	return err
}

func New(cost int) BCrypt {
	return BCrypt{
		cost: cost,
	}
}

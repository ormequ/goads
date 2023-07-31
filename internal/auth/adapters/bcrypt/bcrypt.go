package bcrypt

import (
	"context"
	"errors"
	"goads/internal/auth/app"
	"goads/internal/pkg/errwrap"

	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	cost int
}

func (b BCrypt) Generate(ctx context.Context, password string) (string, error) {
	const op = "bcrypt.Generate"

	if ctx.Err() != nil {
		return "", errwrap.New(ctx.Err(), app.ServiceName, op)
	}
	if len(password) == 0 {
		return "", errwrap.New(app.ErrPasswordToShort, app.ServiceName, op)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		err = errwrap.New(err, app.ServiceName, op)
	}
	return string(bytes), err
}

func (b BCrypt) Compare(ctx context.Context, hash string, password string) error {
	const op = "bcrypt.Compare"

	if ctx.Err() != nil {
		return errwrap.New(ctx.Err(), app.ServiceName, op)
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = errwrap.New(app.ErrIncorrectCredentials, app.ServiceName, op).WithDetails(err.Error())
	} else if err != nil {
		err = errwrap.New(err, app.ServiceName, op)
	}
	return err
}

func New(cost int) BCrypt {
	return BCrypt{
		cost: cost,
	}
}

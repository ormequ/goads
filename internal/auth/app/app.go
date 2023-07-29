package app

import (
	"context"
	"errors"
	structValidator "github.com/ormequ/validator"
	"goads/internal/auth/users"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	Store(ctx context.Context, user users.User) (int64, error)
	GetByID(ctx context.Context, id int64) (users.User, error)
	GetByEmail(ctx context.Context, email string) (users.User, error)
	Update(ctx context.Context, user users.User) error
	Delete(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Tokenizer
type Tokenizer interface {
	Generate(ctx context.Context, id int64) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Hasher
type Hasher interface {
	Generate(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, hash string, password string) error
}

type Validator interface {
	Validate(ctx context.Context, token string) (int64, error)
}

type App struct {
	Repo      Repository
	Tokenizer Tokenizer
	Hasher    Hasher
	Validator Validator
}

func (a App) Register(ctx context.Context, email string, name string, password string) (users.User, error) {
	password, err := a.Hasher.Generate(ctx, password)
	if err != nil {
		return users.User{}, err
	}
	user := users.New(email, name, password)
	err = structValidator.Validate(user)
	if err != nil {
		return users.User{}, errors.Join(err, ErrInvalidContent)
	}

	user.ID, err = a.Repo.Store(ctx, user)
	return user, err
}

func (a App) Authenticate(ctx context.Context, email string, password string) (string, error) {
	user, err := a.Repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = a.Hasher.Compare(ctx, user.Password, password)
	if err != nil {
		return "", err
	}
	return a.Tokenizer.Generate(ctx, user.ID)
}

func (a App) Validate(ctx context.Context, token string) (int64, error) {
	return a.Validator.Validate(ctx, token)
}

func (a App) GetByID(ctx context.Context, id int64) (users.User, error) {
	return a.Repo.GetByID(ctx, id)
}

// change applies function changer for user and updates it in the repository
func (a App) change(ctx context.Context, id int64, changer func(users.User) users.User) (users.User, error) {
	user, err := a.GetByID(ctx, id)
	if err != nil {
		return user, err
	}
	newUser := changer(user)
	err = structValidator.Validate(newUser)
	if err != nil {
		return newUser, errors.Join(err, ErrInvalidContent)
	}
	return newUser, a.Repo.Update(ctx, newUser)
}

func (a App) ChangeEmail(ctx context.Context, id int64, email string) (users.User, error) {
	return a.change(ctx, id, func(user users.User) users.User {
		user.Email = email
		return user
	})
}

func (a App) ChangeName(ctx context.Context, id int64, name string) (users.User, error) {
	return a.change(ctx, id, func(user users.User) users.User {
		user.Name = name
		return user
	})
}

func (a App) ChangePassword(ctx context.Context, id int64, password string) (users.User, error) {
	hash, err := a.Hasher.Generate(ctx, password)
	if err != nil {
		return users.User{}, err
	}
	return a.change(ctx, id, func(user users.User) users.User {
		user.Password = hash
		return user
	})
}

func (a App) Delete(ctx context.Context, id int64) error {
	err := a.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func New(repository Repository, tokenizer Tokenizer, hasher Hasher, validator Validator) App {
	return App{
		Repo:      repository,
		Tokenizer: tokenizer,
		Hasher:    hasher,
		Validator: validator,
	}
}

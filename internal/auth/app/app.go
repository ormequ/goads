package app

import (
	"context"
	"errors"
	"github.com/ormequ/validator"
	"goads/internal/auth/users"
	"goads/internal/pkg/errwrap"
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
	const op = "app.Register"

	password, err := a.Hasher.Generate(ctx, password)
	if err != nil {
		return users.User{}, err
	}
	user := users.New(email, name, password)
	err = govalid.Validate(user)
	if err != nil {
		err = errors.Join(ErrInvalidContent, err)
		return user, errwrap.New(err, ServiceName, op)
	}
	user.ID, err = a.Repo.Store(ctx, user)
	return user, errwrap.JoinWithCaller(err, op)
}

func (a App) Authenticate(ctx context.Context, email string, password string) (string, error) {
	const op = "app.Authenticate"

	user, err := a.Repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errwrap.JoinWithCaller(err, op)
	}
	err = a.Hasher.Compare(ctx, user.Password, password)
	if err != nil {
		return "", errwrap.JoinWithCaller(err, op)
	}
	token, err := a.Tokenizer.Generate(ctx, user.ID)
	return token, errwrap.JoinWithCaller(err, op)
}

func (a App) Validate(ctx context.Context, token string) (int64, error) {
	const op = "app.Validate"
	id, err := a.Validator.Validate(ctx, token)
	return id, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByID(ctx context.Context, id int64) (users.User, error) {
	const op = "app.GetByID"
	user, err := a.Repo.GetByID(ctx, id)
	return user, errwrap.JoinWithCaller(err, op)
}

// change applies function changer for user and updates it in the repository
func (a App) change(ctx context.Context, id int64, changer func(users.User) users.User) (users.User, error) {
	const op = "app.change"
	user, err := a.GetByID(ctx, id)
	if err != nil {
		return user, errwrap.JoinWithCaller(err, op)
	}
	newUser := changer(user)
	err = govalid.Validate(newUser)
	if err != nil {
		err = errors.Join(ErrInvalidContent, err)
		return user, errwrap.New(err, ServiceName, op)
	}
	err = a.Repo.Update(ctx, newUser)
	if err != nil {
		newUser = user
	}
	return newUser, errwrap.JoinWithCaller(err, op)
}

func (a App) ChangeEmail(ctx context.Context, id int64, email string) (users.User, error) {
	const op = "app.ChangeEmail"
	user, err := a.change(ctx, id, func(user users.User) users.User {
		user.Email = email
		return user
	})
	return user, errwrap.JoinWithCaller(err, op)
}

func (a App) ChangeName(ctx context.Context, id int64, name string) (users.User, error) {
	const op = "app.ChangeName"
	user, err := a.change(ctx, id, func(user users.User) users.User {
		user.Name = name
		return user
	})
	return user, errwrap.JoinWithCaller(err, op)
}

func (a App) ChangePassword(ctx context.Context, id int64, password string) (users.User, error) {
	const op = "app.ChangePassword"
	hash, err := a.Hasher.Generate(ctx, password)
	if err != nil {
		return users.User{}, errwrap.JoinWithCaller(err, op)
	}
	user, err := a.change(ctx, id, func(user users.User) users.User {
		user.Password = hash
		return user
	})
	return user, errwrap.JoinWithCaller(err, op)
}

func (a App) Delete(ctx context.Context, id int64) error {
	const op = "app.Delete"
	err := a.Repo.Delete(ctx, id)
	return errwrap.JoinWithCaller(err, op)
}

func New(repository Repository, tokenizer Tokenizer, hasher Hasher, validator Validator) App {
	return App{
		Repo:      repository,
		Tokenizer: tokenizer,
		Hasher:    hasher,
		Validator: validator,
	}
}

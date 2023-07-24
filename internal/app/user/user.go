package user

import (
	"context"
	"fmt"
	validator "github.com/ormequ/validator"
	"goads/internal/app"
	"goads/internal/entities/users"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	Store(ctx context.Context, user users.User) (int64, error)
	GetByID(ctx context.Context, id int64) (users.User, error)
	GetByEmail(ctx context.Context, email string) (users.User, error)
	Update(ctx context.Context, user users.User) error
	Delete(ctx context.Context, id int64) error
}

type App struct {
	repository Repository
}

func (u App) Create(ctx context.Context, email string, name string) (users.User, error) {
	user := users.New(email, name)
	validateErr := validator.Validate(user)
	if validateErr != nil {
		return user, app.Error{
			Err:     app.ErrInvalidContent,
			Details: validateErr.Error(),
		}
	}
	var err error
	user.ID, err = u.repository.Store(ctx, user)
	if err != nil {
		return user, app.Error{
			Err:     err,
			Type:    "user",
			Details: fmt.Sprintf("email: %s, name: %s", email, name),
		}
	}
	return user, nil
}

func (u App) GetByID(ctx context.Context, id int64) (users.User, error) {
	user, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return users.User{}, app.Error{
			Err:  err,
			Type: "user",
			ID:   user.ID,
		}
	}
	return user, nil
}

func (u App) GetByEmail(ctx context.Context, email string) (users.User, error) {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return users.User{}, app.Error{
			Err:  err,
			Type: "user",
			ID:   user.ID,
		}
	}
	return user, nil
}

// change applies function changer for user and updates it in the repository
func (u App) change(ctx context.Context, user users.User, changer func(users.User) users.User) error {
	newUser := changer(user)
	err := validator.Validate(newUser)
	if err != nil {
		return app.Error{
			Err:     app.ErrInvalidContent,
			Type:    "user",
			ID:      user.ID,
			Details: err.Error(),
		}
	}
	err = u.repository.Update(ctx, newUser)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "user",
			ID:   user.ID,
		}
	}
	return nil
}

func (u App) ChangeEmail(ctx context.Context, id int64, email string) error {
	user, err := u.GetByID(ctx, id)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "users",
		}
	}
	return u.change(ctx, user, func(user users.User) users.User {
		user.Email = email
		return user
	})
}

func (u App) ChangeName(ctx context.Context, id int64, name string) error {
	user, err := u.GetByID(ctx, id)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "users",
		}
	}
	return u.change(ctx, user, func(user users.User) users.User {
		user.Name = name
		return user
	})
}

func (u App) Delete(ctx context.Context, id int64) error {
	err := u.repository.Delete(ctx, id)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "users",
		}
	}
	return nil
}

func New(repo Repository) App {
	return App{repository: repo}
}

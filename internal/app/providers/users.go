package providers

import (
	"context"
	"fmt"
	validator "github.com/ormequ/validator"
	"goads/internal/app"
	"goads/internal/entities/users"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=UsersRepository
type UsersRepository interface {
	GetNewID(ctx context.Context) (int64, error)
	Store(ctx context.Context, user users.User) error
	GetByID(ctx context.Context, id int64) (users.User, error)
	Update(ctx context.Context, user users.User) error
	Delete(ctx context.Context, id int64) error
}

type Users struct {
	repository UsersRepository
}

func (u Users) Create(ctx context.Context, email string, name string) (users.User, error) {
	id, err := u.repository.GetNewID(ctx)
	if err != nil {
		return users.User{}, app.Error{
			Err: err,
		}
	}
	user := users.New(id, email, name)
	validateErr := validator.Validate(user)
	if validateErr != nil {
		return users.User{}, app.Error{
			Err:     app.ErrInvalidContent,
			ID:      id,
			Details: validateErr.Error(),
		}
	}
	err = u.repository.Store(ctx, user)
	if err != nil {
		return users.User{}, app.Error{
			Err:  err,
			Type: "user",
			ID:   id,
			Details: fmt.Sprintf(
				"email: %s, name: %s",
				email,
				name,
			),
		}
	}
	return user, nil
}

func (u Users) GetByID(ctx context.Context, id int64) (users.User, error) {
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

// change applies function changer for user and updates it in the repository
func (u Users) change(ctx context.Context, user users.User, changer func(users.User) users.User) error {
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

func (u Users) ChangeEmail(ctx context.Context, id int64, email string) error {
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

func (u Users) ChangeName(ctx context.Context, id int64, name string) error {
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

func (u Users) Delete(ctx context.Context, id int64) error {
	err := u.repository.Delete(ctx, id)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "users",
		}
	}
	return nil
}

func NewUsers(repo UsersRepository) *Users {
	return &Users{repository: repo}
}

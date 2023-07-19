package services

import (
	"context"
	"goads/internal/entities/users"
	"google.golang.org/protobuf/types/known/emptypb"
)

type appUsers interface {
	Create(ctx context.Context, email string, name string) (users.User, error)
	GetByID(ctx context.Context, id int64) (users.User, error)
	GetByEmail(ctx context.Context, email string) (users.User, error)
	ChangeEmail(ctx context.Context, id int64, email string) error
	ChangeName(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error
}

type Users struct {
	app appUsers
}

func (a Users) Create(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := a.app.Create(ctx, request.Email, request.Name)
	return userToResponse(user), GetErrorStatus(err)
}

func (a Users) ChangeName(ctx context.Context, request *ChangeUserNameRequest) (*UserResponse, error) {
	err := a.app.ChangeName(ctx, request.Id, request.Name)
	if err != nil {
		return nil, GetErrorStatus(err)
	}
	user, err := a.app.GetByID(ctx, request.Id)
	return userToResponse(user), GetErrorStatus(err)
}

func (a Users) ChangeEmail(ctx context.Context, request *ChangeUserEmailRequest) (*UserResponse, error) {
	err := a.app.ChangeEmail(ctx, request.Id, request.Email)
	if err != nil {
		return nil, GetErrorStatus(err)
	}
	user, err := a.app.GetByID(ctx, request.Id)
	return userToResponse(user), GetErrorStatus(err)
}

func (a Users) GetByID(ctx context.Context, request *GetUserByIDRequest) (*UserResponse, error) {
	user, err := a.app.GetByID(ctx, request.Id)
	return userToResponse(user), GetErrorStatus(err)
}

func (a Users) GetByEmail(ctx context.Context, request *GetUserByEmailRequest) (*UserResponse, error) {
	user, err := a.app.GetByEmail(ctx, request.Email)
	return userToResponse(user), GetErrorStatus(err)
}

func (a Users) Delete(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	err := a.app.Delete(ctx, request.Id)
	return new(emptypb.Empty), GetErrorStatus(err)
}

func userToResponse(user users.User) *UserResponse {
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func NewUsers(app appUsers) Users {
	return Users{app}
}

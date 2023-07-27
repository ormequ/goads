package grpc

import (
	"context"
	"goads/internal/auth/proto"
	"goads/internal/auth/users"
	"google.golang.org/protobuf/types/known/emptypb"
)

type App interface {
	Register(ctx context.Context, email string, name string, password string) (users.User, error)
	Authenticate(ctx context.Context, email string, password string) (string, error)
	GetByID(ctx context.Context, id int64) (users.User, error)
	ChangeEmail(ctx context.Context, id int64, email string) error
	ChangeName(ctx context.Context, id int64, name string) error
	ChangePassword(ctx context.Context, id int64, password string) error
	Delete(ctx context.Context, id int64) error
	Validate(ctx context.Context, token string) (users.User, error)
}

type Service struct {
	app App
}

func (s Service) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user, err := s.app.Register(ctx, request.Email, request.Name, request.Password)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	token, err := s.app.Authenticate(ctx, request.Email, request.Password)
	return &proto.RegisterResponse{
		User:  userToInfoResponse(user),
		Token: tokenToResponse(token),
	}, getErrorStatus(err)
}

func (s Service) Authenticate(ctx context.Context, request *proto.AuthenticateRequest) (*proto.TokenResponse, error) {
	token, err := s.app.Authenticate(ctx, request.Email, request.Password)
	return tokenToResponse(token), getErrorStatus(err)
}

func (s Service) ChangeName(ctx context.Context, request *proto.ChangeUserNameRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.Validate(ctx, request.Token)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	user.Name = request.Name
	err = s.app.ChangeName(ctx, user.ID, request.Name)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) ChangeEmail(ctx context.Context, request *proto.ChangeUserEmailRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.Validate(ctx, request.Token)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	user.Email = request.Email
	err = s.app.ChangeEmail(ctx, user.ID, request.Email)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) ChangePassword(ctx context.Context, request *proto.ChangeUserPasswordRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.Validate(ctx, request.Token)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	err = s.app.ChangePassword(ctx, user.ID, request.Password)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) GetByID(ctx context.Context, request *proto.GetUserByIDRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.GetByID(ctx, request.Id)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) Delete(ctx context.Context, request *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	user, err := s.app.Validate(ctx, request.Token)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	err = s.app.Delete(ctx, user.ID)
	return new(emptypb.Empty), getErrorStatus(err)
}

func NewService(app App) Service {
	return Service{app}
}

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
	ChangeEmail(ctx context.Context, id int64, email string) (users.User, error)
	ChangeName(ctx context.Context, id int64, name string) (users.User, error)
	ChangePassword(ctx context.Context, id int64, password string) (users.User, error)
	Delete(ctx context.Context, id int64) error
	Validate(ctx context.Context, token string) (int64, error)
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

func (s Service) Validate(ctx context.Context, request *proto.ValidateRequest) (*proto.UserIDResponse, error) {
	id, err := s.app.Validate(ctx, request.Token)
	return &proto.UserIDResponse{Id: id}, getErrorStatus(err)
}

func (s Service) ChangeName(ctx context.Context, request *proto.ChangeUserNameRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.ChangeName(ctx, request.Id, request.Name)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) ChangeEmail(ctx context.Context, request *proto.ChangeUserEmailRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.ChangeEmail(ctx, request.Id, request.Email)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) ChangePassword(ctx context.Context, request *proto.ChangeUserPasswordRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.ChangePassword(ctx, request.Id, request.Password)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) GetByID(ctx context.Context, request *proto.GetUserByIDRequest) (*proto.UserInfoResponse, error) {
	user, err := s.app.GetByID(ctx, request.Id)
	return userToInfoResponse(user), getErrorStatus(err)
}

func (s Service) Delete(ctx context.Context, request *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), getErrorStatus(s.app.Delete(ctx, request.Id))
}

func NewService(app App) Service {
	return Service{app}
}

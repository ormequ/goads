package grpc

import (
	"errors"
	"goads/internal/auth/app"
	"goads/internal/auth/proto"
	"goads/internal/auth/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorStatus(err error) error {
	if err == nil {
		return nil
	}
	code := codes.Internal
	if errors.Is(err, app.ErrNotFound) {
		code = codes.NotFound
	}
	if errors.Is(err, app.ErrIncorrectCredentials) || errors.Is(err, app.ErrInvalidToken) {
		code = codes.Unauthenticated
	}
	if errors.Is(err, app.ErrInvalidContent) {
		code = codes.InvalidArgument
	}
	if errors.Is(err, app.ErrEmailAlreadyExists) {
		code = codes.AlreadyExists
	}
	return status.Error(code, err.Error())
}

func userToInfoResponse(user users.User) *proto.UserInfoResponse {
	return &proto.UserInfoResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func tokenToResponse(token string) *proto.TokenResponse {
	return &proto.TokenResponse{
		Token: token,
	}
}

package grpc

import (
	"context"
	"errors"
	"goads/internal/app"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"goads/internal/filters"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Service struct {
	Ads   app.Ads
	Users app.Users
}

func (a Service) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := a.Ads.Create(ctx, request.Title, request.Text, request.UserId)
	return adToResponse(ad), getErrorStatus(err)
}

func (a Service) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	err := a.Ads.ChangeStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	ad, err := a.Ads.GetByID(ctx, request.AdId)
	return adToResponse(ad), getErrorStatus(err)
}

func (a Service) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	err := a.Ads.Update(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	ad, err := a.Ads.GetByID(ctx, request.AdId)
	return adToResponse(ad), getErrorStatus(err)
}

func (a Service) ListAds(ctx context.Context, request *FilterAdsRequest) (*ListAdResponse, error) {
	list, err := a.Ads.GetFiltered(ctx, filters.AdsOptions{
		AuthorID: request.AuthorId,
		Date:     time.UnixMilli(request.Date).UTC(),
		All:      request.All,
	})
	return adsListToResponse(list), getErrorStatus(err)
}

func (a Service) SearchAds(ctx context.Context, request *SearchAdsRequest) (*ListAdResponse, error) {
	list, err := a.Ads.Search(ctx, request.Title)
	return adsListToResponse(list), getErrorStatus(err)
}

func (a Service) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	err := a.Ads.Delete(ctx, request.AdId, request.AuthorId)
	return new(emptypb.Empty), getErrorStatus(err)
}

func (a Service) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := a.Users.Create(ctx, request.Email, request.Name)
	return userToResponse(user), getErrorStatus(err)
}

func (a Service) ChangeUserName(ctx context.Context, request *ChangeUserNameRequest) (*UserResponse, error) {
	err := a.Users.ChangeName(ctx, request.Id, request.Name)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	user, err := a.Users.GetByID(ctx, request.Id)
	return userToResponse(user), getErrorStatus(err)
}

func (a Service) ChangeUserEmail(ctx context.Context, request *ChangeUserEmailRequest) (*UserResponse, error) {
	err := a.Users.ChangeEmail(ctx, request.Id, request.Email)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	user, err := a.Users.GetByID(ctx, request.Id)
	return userToResponse(user), getErrorStatus(err)
}

func (a Service) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	user, err := a.Users.GetByID(ctx, request.Id)
	return userToResponse(user), getErrorStatus(err)
}

func (a Service) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	err := a.Users.Delete(ctx, request.Id)
	return new(emptypb.Empty), getErrorStatus(err)
}

func adToResponse(ad ads.Ad) *AdResponse {
	return &AdResponse{
		Id:         ad.ID,
		Title:      ad.Title,
		Published:  ad.Published,
		Text:       ad.Text,
		AuthorId:   ad.AuthorID,
		CreateDate: ad.CreateDate.UnixMilli(),
		UpdateDate: ad.UpdateDate.UnixMilli(),
	}
}

func adsListToResponse(list []ads.Ad) *ListAdResponse {
	res := ListAdResponse{
		List: make([]*AdResponse, len(list)),
	}
	for i := range list {
		res.List[i] = adToResponse(list[i])
	}
	return &res
}

func userToResponse(user users.User) *UserResponse {
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func NewService(a app.Ads, u app.Users) Service {
	return Service{a, u}
}

func getErrorStatus(err error) error {
	if err == nil {
		return nil
	}
	code := codes.Internal
	if errors.Is(err, app.ErrNotFound) {
		code = codes.NotFound
	}
	if errors.Is(err, app.ErrPermissionDenied) {
		code = codes.PermissionDenied
	}
	if errors.Is(err, app.ErrInvalidContent) || errors.Is(err, app.ErrInvalidFilter) {
		code = codes.InvalidArgument
	}
	return status.Error(code, err.Error())
}

package grpc

import (
	"context"
	"goads/internal/ads/ads"
	"goads/internal/ads/app"
	"goads/internal/ads/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type App interface {
	Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error)
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
	ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error
	Update(ctx context.Context, id int64, userID int64, title string, text string) error
	GetFiltered(ctx context.Context, opt app.Filter) ([]ads.Ad, error)
	Delete(ctx context.Context, id int64, userID int64) error
	Search(ctx context.Context, title string) ([]ads.Ad, error)
}

type Service struct {
	app App
}

func (a Service) Create(ctx context.Context, request *proto.CreateAdRequest) (*proto.AdResponse, error) {
	ad, err := a.app.Create(ctx, request.Title, request.Text, request.UserId)
	return adToResponse(ad), getErrorStatus(err)
}

func (a Service) ChangeStatus(ctx context.Context, request *proto.ChangeAdStatusRequest) (*proto.AdResponse, error) {
	err := a.app.ChangeStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	ad, err := a.app.GetByID(ctx, request.AdId)
	return adToResponse(ad), getErrorStatus(err)
}

func (a Service) Update(ctx context.Context, request *proto.UpdateAdRequest) (*proto.AdResponse, error) {
	err := a.app.Update(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	ad, err := a.app.GetByID(ctx, request.AdId)
	return adToResponse(ad), getErrorStatus(err)
}

func (a Service) List(ctx context.Context, request *proto.FilterAdsRequest) (*proto.ListAdResponse, error) {
	list, err := a.app.GetFiltered(ctx, app.Filter{
		AuthorID: request.AuthorId,
		Date:     time.UnixMilli(request.Date).UTC(),
		All:      request.All,
	})
	return adsListToResponse(list), getErrorStatus(err)
}

func (a Service) Search(ctx context.Context, request *proto.SearchAdsRequest) (*proto.ListAdResponse, error) {
	list, err := a.app.Search(ctx, request.Title)
	return adsListToResponse(list), getErrorStatus(err)
}

func (a Service) Delete(ctx context.Context, request *proto.DeleteAdRequest) (*emptypb.Empty, error) {
	err := a.app.Delete(ctx, request.AdId, request.AuthorId)
	return new(emptypb.Empty), getErrorStatus(err)
}

func NewService(app App) Service {
	return Service{app}
}

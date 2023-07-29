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
	ChangeStatus(ctx context.Context, id int64, userID int64, published bool) (ads.Ad, error)
	Update(ctx context.Context, id int64, userID int64, title string, text string) (ads.Ad, error)
	GetFiltered(ctx context.Context, opt app.Filter) ([]ads.Ad, error)
	Delete(ctx context.Context, id int64, userID int64) error
	Search(ctx context.Context, title string) ([]ads.Ad, error)
	GetOnlyPublished(ctx context.Context, ids []int64) ([]ads.Ad, error)
}

type Service struct {
	app App
}

func (s Service) GetOnlyPublished(ctx context.Context, request *proto.AdIDsRequest) (*proto.ListAdResponse, error) {
	list, err := s.app.GetOnlyPublished(ctx, request.Id)
	return adsListToResponse(list), getErrorStatus(err)
}

func (s Service) Create(ctx context.Context, request *proto.CreateAdRequest) (*proto.AdResponse, error) {
	ad, err := s.app.Create(ctx, request.Title, request.Text, request.AuthorId)
	return adToResponse(ad), getErrorStatus(err)
}

func (s Service) ChangeStatus(ctx context.Context, request *proto.ChangeAdStatusRequest) (*proto.AdResponse, error) {
	ad, err := s.app.ChangeStatus(ctx, request.AdId, request.AuthorId, request.Published)
	return adToResponse(ad), getErrorStatus(err)
}

func (s Service) Update(ctx context.Context, request *proto.UpdateAdRequest) (*proto.AdResponse, error) {
	ad, err := s.app.Update(ctx, request.AdId, request.AuthorId, request.Title, request.Text)
	return adToResponse(ad), getErrorStatus(err)
}

func (s Service) Filter(ctx context.Context, request *proto.FilterAdsRequest) (*proto.ListAdResponse, error) {
	list, err := s.app.GetFiltered(ctx, app.Filter{
		AuthorID: request.AuthorId,
		Date:     time.UnixMilli(request.Date).UTC(),
		All:      request.All,
		Title:    request.Title,
	})
	return adsListToResponse(list), getErrorStatus(err)
}

func (s Service) GetByID(ctx context.Context, request *proto.GetAdByIDRequest) (*proto.AdResponse, error) {
	ad, err := s.app.GetByID(ctx, request.Id)
	return adToResponse(ad), getErrorStatus(err)
}

func (s Service) Delete(ctx context.Context, request *proto.DeleteAdRequest) (*emptypb.Empty, error) {
	err := s.app.Delete(ctx, request.AdId, request.AuthorId)
	return new(emptypb.Empty), getErrorStatus(err)
}

func NewService(app App) Service {
	return Service{app}
}

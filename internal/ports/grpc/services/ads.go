package services

import (
	"context"
	adApp "goads/internal/app/ad"
	"goads/internal/entities/ads"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type appAds interface {
	Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error)
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
	ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error
	Update(ctx context.Context, id int64, userID int64, title string, text string) error
	GetFiltered(ctx context.Context, opt adApp.Filter) ([]ads.Ad, error)
	Delete(ctx context.Context, id int64, userID int64) error
	Search(ctx context.Context, title string) ([]ads.Ad, error)
}

type Ads struct {
	app appAds
}

func (a Ads) Create(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := a.app.Create(ctx, request.Title, request.Text, request.UserId)
	return adToResponse(ad), GetErrorStatus(err)
}

func (a Ads) ChangeStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	err := a.app.ChangeStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, GetErrorStatus(err)
	}
	ad, err := a.app.GetByID(ctx, request.AdId)
	return adToResponse(ad), GetErrorStatus(err)
}

func (a Ads) Update(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	err := a.app.Update(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, GetErrorStatus(err)
	}
	ad, err := a.app.GetByID(ctx, request.AdId)
	return adToResponse(ad), GetErrorStatus(err)
}

func (a Ads) List(ctx context.Context, request *FilterAdsRequest) (*ListAdResponse, error) {
	list, err := a.app.GetFiltered(ctx, adApp.Filter{
		AuthorID: request.AuthorId,
		Date:     time.UnixMilli(request.Date).UTC(),
		All:      request.All,
	})
	return adsListToResponse(list), GetErrorStatus(err)
}

func (a Ads) Search(ctx context.Context, request *SearchAdsRequest) (*ListAdResponse, error) {
	list, err := a.app.Search(ctx, request.Title)
	return adsListToResponse(list), GetErrorStatus(err)
}

func (a Ads) Delete(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	err := a.app.Delete(ctx, request.AdId, request.AuthorId)
	return new(emptypb.Empty), GetErrorStatus(err)
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

func NewAds(app appAds) Ads {
	return Ads{app}
}

package maprepo

import (
	"context"
	"goads/internal/app"
	adApp "goads/internal/app/ad"
	"goads/internal/entities/ads"
	"strings"
	"time"
)

type Ads struct {
	Repo
}

func (r *Ads) Store(ctx context.Context, ad ads.Ad) (int64, error) {
	id, err := r.Repo.GetNewID(ctx)
	if err != nil {
		return id, err
	}
	ad.ID = id
	return id, r.Repo.Store(ctx, ad, id)
}

func (r *Ads) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
	ent, err := r.Repo.GetByID(ctx, id)
	if err != nil {
		return ads.Ad{}, err
	}
	ad, ok := ent.(ads.Ad)
	if !ok {
		return ads.Ad{}, app.ErrNotFound
	}
	return ad, nil
}

func (r *Ads) GetFiltered(ctx context.Context, filter adApp.Filter) ([]ads.Ad, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	res := make([]ads.Ad, 0, r.size.Load())
	date := filter.Date.Truncate(time.Hour * 24)
	r.storage.Range(func(key, v any) bool {
		ad := v.(ads.Ad)
		if (ad.Published || filter.All) &&
			(ad.AuthorID == filter.AuthorID || filter.AuthorID == -1) &&
			(ad.GetCreateDate() == date || filter.Date.IsZero() || filter.Date.UnixMilli() == 0) &&
			strings.HasPrefix(ad.Title, filter.Prefix) {

			res = append(res, ad)
		}
		return true
	})
	return res, nil
}

func (r *Ads) Update(ctx context.Context, ad ads.Ad) error {
	return r.Repo.Update(ctx, ad, ad.ID)
}

func NewAds() *Ads {
	return &Ads{}
}

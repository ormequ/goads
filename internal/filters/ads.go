package filters

import (
	"goads/internal/entities"
	"goads/internal/entities/ads"
	"time"
)

type AdsOptions struct {
	AuthorID int64
	Date     time.Time
	All      bool
}

func AdByAuthor(opt AdsOptions) func(entities.Interface) bool {
	return func(v entities.Interface) bool {
		ad, ok := v.(ads.Ad)
		return ok && (ad.AuthorID == opt.AuthorID || opt.AuthorID == -1)
	}
}

func AdByCreateDate(opt AdsOptions) func(entities.Interface) bool {
	return func(v entities.Interface) bool {
		ad, ok := v.(ads.Ad)
		isZero := opt.Date.IsZero() || opt.Date.UnixMilli() == 0
		return ok && (ad.CreateDate.Truncate(time.Second) == opt.Date.Truncate(time.Second) || isZero)
	}
}

func AdPublished(opt AdsOptions) func(p entities.Interface) bool {
	return func(v entities.Interface) bool {
		ad, ok := v.(ads.Ad)
		return ok && (ad.Published || opt.All)
	}
}

func NewAds(opt AdsOptions) entities.Filter {
	return entities.Filter{
		AdByAuthor(opt),
		AdByCreateDate(opt),
		AdPublished(opt),
	}
}

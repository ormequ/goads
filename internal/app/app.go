package app

import (
	"context"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"goads/internal/filters"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Ads
type Ads interface {
	Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error)
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
	ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error
	Update(ctx context.Context, id int64, userID int64, title string, text string) error
	GetFiltered(ctx context.Context, filter filters.AdsOptions) ([]ads.Ad, error)
	Search(ctx context.Context, title string) ([]ads.Ad, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Users
type Users interface {
	Create(ctx context.Context, email string, name string) (users.User, error)
	GetByID(ctx context.Context, id int64) (users.User, error)
	ChangeEmail(ctx context.Context, id int64, email string) error
	ChangeName(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error
}

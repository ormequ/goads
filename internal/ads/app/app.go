package app

import (
	"context"
	"errors"
	"fmt"
	validator "github.com/ormequ/validator"
	"goads/internal/ads/ads"
	"goads/internal/pkg/errwrap"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	Store(ctx context.Context, ad ads.Ad) (int64, error)
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
	GetFiltered(ctx context.Context, filter Filter) ([]ads.Ad, error)
	Update(ctx context.Context, ad ads.Ad) error
	Delete(ctx context.Context, id int64) error
	GetOnlyPublished(ctx context.Context, ids []int64) ([]ads.Ad, error)
}

type Filter struct {
	AuthorID int64
	Date     time.Time
	Title    string
	All      bool
}

type App struct {
	repository Repository
}

// Create creates a new ad with incremented id
func (a App) Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error) {
	const op = "app.Create"

	ad := ads.New(title, text, authorID)
	err := validator.Validate(ad)
	if err != nil {
		err = errors.Join(ErrInvalidContent, err)
		return ad, errwrap.New(err, ServiceName, op)
	}
	ad.ID, err = a.repository.Store(ctx, ad)
	return ad, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
	const op = "app.GetByID"
	ad, err := a.repository.GetByID(ctx, id)
	return ad, errwrap.JoinWithCaller(err, op)
}

func (a App) getEditable(ctx context.Context, id int64, userID int64) (ads.Ad, error) {
	const op = "app.getEditable"

	ad, err := a.repository.GetByID(ctx, id)
	if err != nil {
		return ad, errwrap.JoinWithCaller(err, op)
	}
	if ad.AuthorID != userID {
		return ad, errwrap.New(ErrPermissionDenied, ServiceName, op).OnObject("ad", ad.ID).
			WithDetails(fmt.Sprintf("ad created by %d and cannot be changed by %d", ad.AuthorID, userID))
	}
	return ad, nil
}

// change applies function changer for an ad and updates it in the repository
func (a App) change(ctx context.Context, id int64, authorID int64, changer func(ads.Ad) ads.Ad) (ads.Ad, error) {
	const op = "app.change"

	ad, err := a.getEditable(ctx, id, authorID)
	if err != nil {
		return ad, errwrap.JoinWithCaller(err, op)
	}
	newAd := changer(ad)
	err = validator.Validate(newAd)
	if err != nil {
		err = errors.Join(ErrInvalidContent, err)
		return ad, errwrap.New(err, ServiceName, op)
	}
	newAd.UpdateDate = time.Now().UTC()
	err = a.repository.Update(ctx, newAd)
	return ad, errwrap.JoinWithCaller(err, op)
}

// ChangeStatus changes ad's status only if userID is equal to author id of the ad
func (a App) ChangeStatus(ctx context.Context, id int64, authorID int64, published bool) (ads.Ad, error) {
	const op = "app.ChangeStatus"
	ad, err := a.change(ctx, id, authorID, func(ad ads.Ad) ads.Ad {
		ad.Published = published
		return ad
	})
	return ad, errwrap.JoinWithCaller(err, op)
}

// Update changes ad's content (title and text) only if userID is equal to author id of the ad
func (a App) Update(ctx context.Context, id int64, authorID int64, title string, text string) (ads.Ad, error) {
	const op = "app.Update"
	ad, err := a.change(ctx, id, authorID, func(ad ads.Ad) ads.Ad {
		ad.Title, ad.Text = title, text
		return ad
	})
	return ad, errwrap.JoinWithCaller(err, op)
}

// GetFiltered returns list of ads satisfy filter. Usable filters look at /internal/filters
func (a App) GetFiltered(ctx context.Context, opt Filter) ([]ads.Ad, error) {
	const op = "app.GetFiltered"
	list, err := a.repository.GetFiltered(ctx, opt)
	return list, errwrap.JoinWithCaller(err, op)
}

func (a App) GetOnlyPublished(ctx context.Context, ids []int64) ([]ads.Ad, error) {
	const op = "app.GetOnlyPublished"
	list, err := a.repository.GetOnlyPublished(ctx, ids)
	return list, errwrap.JoinWithCaller(err, op)
}

// Delete removes ad with got id if userID equals to author ID of the ad
func (a App) Delete(ctx context.Context, id int64, userID int64) error {
	const op = "app.Delete"
	_, err := a.getEditable(ctx, id, userID)
	if err != nil {
		return errwrap.JoinWithCaller(err, op)
	}
	err = a.repository.Delete(ctx, id)
	return errwrap.JoinWithCaller(err, op)
}

// Search finds ads with prefix equals to title
func (a App) Search(ctx context.Context, title string) ([]ads.Ad, error) {
	const op = "app.Search"
	list, err := a.repository.GetFiltered(ctx, Filter{Title: title, AuthorID: -1, All: true})
	return list, errwrap.JoinWithCaller(err, op)
}

func New(repo Repository) App {
	return App{repository: repo}
}

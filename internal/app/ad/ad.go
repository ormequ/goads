package ad

import (
	"context"
	"fmt"
	validator "github.com/ormequ/validator"
	"goads/internal/app"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	Store(ctx context.Context, ad ads.Ad) (int64, error)
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
	GetFiltered(ctx context.Context, filter Filter) ([]ads.Ad, error)
	Update(ctx context.Context, ad ads.Ad) error
	Delete(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=UsersGetter
type UsersGetter interface {
	GetByID(ctx context.Context, id int64) (users.User, error)
}

type Filter struct {
	AuthorID int64
	Date     time.Time
	Prefix   string
	All      bool
}

type App struct {
	repository Repository
	users      UsersGetter
}

// Create creates a new ad with incremented id
func (a App) Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error) {
	_, err := a.users.GetByID(ctx, authorID)
	if err != nil {
		return ads.Ad{}, err
	}
	ad := ads.New(title, text, authorID)
	validateErr := validator.Validate(ad)
	if validateErr != nil {
		return ad, app.Error{
			Err:     app.ErrInvalidContent,
			Type:    "ad",
			Details: validateErr.Error(),
		}
	}
	ad.ID, err = a.repository.Store(ctx, ad)
	if err != nil {
		return ad, app.Error{
			Err:     err,
			Type:    "ad",
			Details: fmt.Sprintf("title: %s, text: %s, author: %d", title, text, authorID),
		}
	}
	return ad, nil
}

func (a App) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
	ad, err := a.repository.GetByID(ctx, id)
	if err != nil {
		return ads.Ad{}, app.Error{
			Err:  err,
			Type: "ad",
			ID:   ad.ID,
		}
	}
	return ad, nil
}

func (a App) getEditable(ctx context.Context, id int64, userID int64) (ads.Ad, error) {
	ad, err := a.repository.GetByID(ctx, id)
	if err != nil {
		return ads.Ad{}, app.Error{
			Err:  err,
			Type: "ad",
			ID:   ad.ID,
		}
	}
	_, err = a.users.GetByID(ctx, userID)
	if err != nil {
		return ads.Ad{}, err
	}
	if ad.AuthorID != userID {
		return ads.Ad{}, app.Error{
			Err:     app.ErrPermissionDenied,
			Type:    "ad",
			ID:      ad.ID,
			Details: fmt.Sprintf("ad created by %d and cannot be changed by %d", ad.AuthorID, userID),
		}
	}
	return ad, nil
}

// change applies function changer for an ad and updates it in the repository
func (a App) change(ctx context.Context, ad ads.Ad, changer func(ads.Ad) ads.Ad) error {
	newAd := changer(ad)
	err := validator.Validate(newAd)
	if err != nil {
		return app.Error{
			Err:     app.ErrInvalidContent,
			Type:    "ad",
			ID:      ad.ID,
			Details: err.Error(),
		}
	}
	newAd.UpdateDate = time.Now().UTC()
	err = a.repository.Update(ctx, newAd)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "ad",
			ID:   ad.ID,
		}
	}
	return nil
}

// ChangeStatus changes ad's status only if userID is equal to author id of the ad
func (a App) ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error {
	ad, err := a.getEditable(ctx, id, userID)
	if err != nil {
		return err
	}
	return a.change(ctx, ad, func(ad ads.Ad) ads.Ad {
		ad.Published = published
		return ad
	})
}

// Update changes ad's content (title and text) only if userID is equal to author id of the ad
func (a App) Update(ctx context.Context, id int64, userID int64, title string, text string) error {
	ad, err := a.getEditable(ctx, id, userID)
	if err != nil {
		return err
	}
	return a.change(ctx, ad, func(ad ads.Ad) ads.Ad {
		ad.Title, ad.Text = title, text
		return ad
	})
}

// GetFiltered returns list of ads satisfy filter. Usable filters look at /internal/filters
func (a App) GetFiltered(ctx context.Context, opt Filter) ([]ads.Ad, error) {
	if opt.AuthorID != -1 {
		_, err := a.users.GetByID(ctx, opt.AuthorID)
		if err != nil {
			return nil, app.Error{
				Err:     app.ErrInvalidFilter,
				Type:    "ads",
				Details: err.Error(),
			}
		}
	}
	list, err := a.repository.GetFiltered(ctx, opt)
	if err != nil {
		return nil, app.Error{
			Err:  err,
			Type: "ads",
		}
	}
	return list, nil
}

// Delete removes ad with got id if userID equals to author ID of the ad
func (a App) Delete(ctx context.Context, id int64, userID int64) error {
	_, err := a.getEditable(ctx, id, userID)
	if err != nil {
		return err
	}
	err = a.repository.Delete(ctx, id)
	if err != nil {
		return app.Error{
			Err:  err,
			Type: "ad",
			ID:   id,
		}
	}
	return nil
}

// Search finds ads with prefix equals to title
func (a App) Search(ctx context.Context, title string) ([]ads.Ad, error) {
	list, err := a.repository.GetFiltered(ctx, Filter{Prefix: title, AuthorID: -1, All: true})
	if err != nil {
		return nil, app.Error{
			Err:  err,
			Type: "ads",
		}
	}
	return list, nil
}

func New(repo Repository, users UsersGetter) App {
	return App{repository: repo, users: users}
}

package providers

import (
	"context"
	"fmt"
	validator "github.com/ormequ/validator"
	"goads/internal/app"
	"goads/internal/entities"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"goads/internal/filters"
	"strings"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=AdsRepository
type AdsRepository interface {
	GetNewID(ctx context.Context) (int64, error)
	Store(ctx context.Context, ad ads.Ad) error
	GetByID(ctx context.Context, id int64) (ads.Ad, error)
	GetFiltered(ctx context.Context, filter entities.Filter) ([]ads.Ad, error)
	Update(ctx context.Context, ad ads.Ad) error
	Delete(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=UsersGetter
type UsersGetter interface {
	GetByID(ctx context.Context, id int64) (users.User, error)
}

type Ads struct {
	repository AdsRepository
	users      UsersGetter
}

// Create creates a new ad with incremented id
func (a Ads) Create(ctx context.Context, title string, text string, authorID int64) (ads.Ad, error) {
	id, err := a.repository.GetNewID(ctx)
	if err != nil {
		return ads.Ad{}, app.Error{
			Err:  err,
			Type: "ad",
		}
	}
	_, err = a.users.GetByID(ctx, authorID)
	if err != nil {
		return ads.Ad{}, err
	}
	ad := ads.New(id, title, text, authorID)
	validateErr := validator.Validate(ad)
	if validateErr != nil {
		return ads.Ad{}, app.Error{
			Err:     app.ErrInvalidContent,
			Type:    "ad",
			ID:      id,
			Details: validateErr.Error(),
		}
	}
	err = a.repository.Store(ctx, ad)
	if err != nil {
		return ads.Ad{}, app.Error{
			Err:  err,
			Type: "ad",
			ID:   id,
			Details: fmt.Sprintf(
				"title: %s, text: %s, author: %d",
				title,
				text,
				authorID,
			),
		}
	}
	return ad, nil
}

func (a Ads) GetByID(ctx context.Context, id int64) (ads.Ad, error) {
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

func (a Ads) getEditable(ctx context.Context, id int64, userID int64) (ads.Ad, error) {
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
			Err:  app.ErrPermissionDenied,
			Type: "ad",
			ID:   ad.ID,
			Details: fmt.Sprintf(
				"ad created by %d and cannot be changed by %d",
				ad.AuthorID,
				userID,
			),
		}
	}
	return ad, nil
}

// change applies function changer for an ad and updates it in the repository
func (a Ads) change(ctx context.Context, ad ads.Ad, changer func(ads.Ad) ads.Ad) error {
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
func (a Ads) ChangeStatus(ctx context.Context, id int64, userID int64, published bool) error {
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
func (a Ads) Update(ctx context.Context, id int64, userID int64, title string, text string) error {
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
func (a Ads) GetFiltered(ctx context.Context, opt filters.AdsOptions) ([]ads.Ad, error) {
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
	list, err := a.repository.GetFiltered(ctx, filters.NewAds(opt))
	if err != nil {
		return nil, app.Error{
			Err:  err,
			Type: "ads",
		}
	}
	return list, nil
}

// Delete removes ad with got id if userID equals to author ID of the ad
func (a Ads) Delete(ctx context.Context, id int64, userID int64) error {
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
func (a Ads) Search(ctx context.Context, title string) ([]ads.Ad, error) {
	list, err := a.repository.GetFiltered(ctx, entities.Filter{func(v entities.Interface) bool {
		ad, ok := v.(ads.Ad)
		return ok && strings.HasPrefix(ad.Title, title)
	}})
	if err != nil {
		return nil, app.Error{
			Err:  err,
			Type: "ads",
		}
	}
	return list, nil
}

func NewAds(repo AdsRepository, users *Users) *Ads {
	return &Ads{repository: repo, users: users}
}

package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ormequ/validator"
	"goads/internal/pkg/errwrap"
	"goads/internal/urlshortener/entities/ads"
	"goads/internal/urlshortener/entities/links"
	"goads/internal/urlshortener/entities/redirects"
	"math/rand"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	Store(ctx context.Context, link links.Link) (int64, error)
	GetByID(ctx context.Context, id int64) (links.Link, error)
	GetByAuthor(ctx context.Context, authorID int64) ([]links.Link, error)
	GetByAlias(ctx context.Context, alias string) (links.Link, error)
	UpdateAlias(ctx context.Context, id int64, alias string) error
	AddAd(ctx context.Context, linkID int64, adID int64) error
	DeleteAd(ctx context.Context, linkID int64, adID int64) error
	Delete(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Generator
type Generator interface {
	Generate(ctx context.Context) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=AdsService
type AdsService interface {
	GetOnlyPublished(ctx context.Context, ids []int64) ([]ads.Ad, error)
}

type App struct {
	Repo Repository
	Gen  Generator
	Ads  AdsService
}

func (a App) generateFreeAlias(ctx context.Context) (alias string, err error) {
	const op = "app.generateFreeAlias"

	var getErr error = nil
	for !errors.Is(getErr, ErrNotFound) && err == nil {
		alias, err = a.Gen.Generate(ctx)
		_, getErr = a.GetByAlias(ctx, alias)
	}
	err = errwrap.JoinWithCaller(err, op)
	return
}

func (a App) Create(ctx context.Context, url string, alias string, authorID int64, ads []int64) (links.Link, error) {
	const op = "app.generateFreeAlias"

	var err error
	if alias == "" {
		alias, err = a.generateFreeAlias(ctx)
		if err != nil {
			return links.Link{}, errwrap.JoinWithCaller(err, op)
		}
	}
	link := links.New(url, alias, authorID, ads)
	err = govalid.Validate(link)
	if err != nil {
		err = errors.Join(ErrInvalidContent, err)
		return link, errwrap.New(err, ServiceName, op)
	}
	link.ID, err = a.Repo.Store(ctx, link)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByID(ctx context.Context, id int64) (links.Link, error) {
	const op = "app.GetByID"
	link, err := a.Repo.GetByID(ctx, id)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByAuthor(ctx context.Context, author int64) ([]links.Link, error) {
	const op = "app.GetByAuthor"
	link, err := a.Repo.GetByAuthor(ctx, author)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) GetRedirect(ctx context.Context, alias string) (redirects.Redirect, error) {
	const op = "app.GetRedirect"
	link, err := a.GetByAlias(ctx, alias)
	if err != nil {
		return redirects.Redirect{}, errwrap.JoinWithCaller(err, op)
	}

	adsList, err := a.Ads.GetOnlyPublished(ctx, link.Ads)
	var ad ads.Ad
	if err == nil && len(adsList) > 0 {
		ad = adsList[rand.Intn(len(adsList))]
	}
	return redirects.New(link, ad), nil
}

func (a App) GetByAlias(ctx context.Context, alias string) (links.Link, error) {
	const op = "app.GetByAliasWithAd"
	link, err := a.Repo.GetByAlias(ctx, alias)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) getEditable(ctx context.Context, id int64, authorID int64) (links.Link, error) {
	const op = "app.getEditable"
	link, err := a.GetByID(ctx, id)
	if err != nil {
		return link, errwrap.JoinWithCaller(err, op)
	}
	if link.AuthorID != authorID {
		return link, errwrap.New(ErrPermissionDenied, ServiceName, op).OnObject("link", id)
	}
	return link, nil
}

func (a App) UpdateAlias(ctx context.Context, id int64, authorID int64, alias string) (links.Link, error) {
	const op = "app.UpdateAlias"
	link, err := a.getEditable(ctx, id, authorID)
	if err != nil {
		return link, errwrap.JoinWithCaller(err, op)
	}
	if alias == "" {
		alias, err = a.generateFreeAlias(ctx)
		if err != nil {
			return link, errwrap.JoinWithCaller(err, op)
		}
	}
	prev := link.Alias
	link.Alias = alias
	err = govalid.Validate(link)
	if err != nil {
		link.Alias = prev
		err = errors.Join(ErrInvalidContent, err)
		return link, errwrap.New(err, ServiceName, op).OnObject("link", id)
	}
	err = a.Repo.UpdateAlias(ctx, id, alias)
	if err != nil {
		link.Alias = prev
	}
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) AddAd(ctx context.Context, linkID int64, adID int64, authorID int64) (links.Link, error) {
	const op = "app.AddAd"
	link, err := a.getEditable(ctx, linkID, authorID)
	if err != nil {
		return link, errwrap.JoinWithCaller(err, op)
	}
	for i := range link.Ads {
		if link.Ads[i] == adID {
			return link, errwrap.New(ErrAdAlreadyAdded, ServiceName, op).
				OnObject("link", linkID).
				WithDetails(fmt.Sprintf("ad ID: %d", adID))
		}
	}
	err = a.Repo.AddAd(ctx, linkID, adID)
	if err == nil {
		link.Ads = append(link.Ads, adID)
	}
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) DeleteAd(ctx context.Context, linkID int64, adID int64, authorID int64) (links.Link, error) {
	const op = "app.DeleteAd"
	link, err := a.getEditable(ctx, linkID, authorID)
	if err != nil {
		return link, errwrap.JoinWithCaller(err, op)
	}
	adIdx := -1
	for i, v := range link.Ads {
		if v == adID {
			adIdx = i
			break
		}
	}
	if adIdx == -1 {
		return link, errwrap.New(ErrNotFound, ServiceName, op).
			OnObject("link", linkID).
			WithDetails(fmt.Sprintf("ad ID: %d", adID))
	}
	err = a.Repo.DeleteAd(ctx, linkID, adID)
	if err == nil {
		link.Ads = append(link.Ads[:adIdx], link.Ads[adIdx+1:]...)
	}
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) Delete(ctx context.Context, id int64, authorID int64) error {
	const op = "app.Delete"
	_, err := a.getEditable(ctx, id, authorID)
	if err != nil {
		return errwrap.JoinWithCaller(err, op)
	}
	return a.Repo.Delete(ctx, id)
}

func New(repo Repository, generator Generator, ads AdsService) App {
	return App{
		Repo: repo,
		Gen:  generator,
		Ads:  ads,
	}
}

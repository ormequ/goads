package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ormequ/validator"
	"goads/internal/pkg/errwrap"
	"goads/internal/urlshortener/links"
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

type App struct {
	repo Repository
	gen  Generator
}

func (a App) generateFreeAlias(ctx context.Context) (alias string, err error) {
	const op = "app.generateFreeAlias"

	var getErr error = nil
	for !errors.Is(getErr, ErrNotFound) && err == nil {
		alias, err = a.gen.Generate(ctx)
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
	link.ID, err = a.repo.Store(ctx, link)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByID(ctx context.Context, id int64) (links.Link, error) {
	const op = "app.GetByID"
	link, err := a.repo.GetByID(ctx, id)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByAuthor(ctx context.Context, author int64) ([]links.Link, error) {
	const op = "app.GetByAuthor"
	link, err := a.repo.GetByAuthor(ctx, author)
	return link, errwrap.JoinWithCaller(err, op)
}

func (a App) GetByAlias(ctx context.Context, alias string) (links.Link, error) {
	const op = "app.GetByAlias"
	link, err := a.repo.GetByAlias(ctx, alias)
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
	err = a.repo.UpdateAlias(ctx, id, alias)
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
	err = a.repo.AddAd(ctx, linkID, adID)
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
	err = a.repo.DeleteAd(ctx, linkID, adID)
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
	return a.repo.Delete(ctx, id)
}

func New(repo Repository, generator Generator) App {
	return App{
		repo: repo,
		gen:  generator,
	}
}

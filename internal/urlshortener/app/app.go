package app

import (
	"context"
	"errors"
	validator "github.com/ormequ/validator"
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
	var getErr error = nil
	for !errors.Is(getErr, ErrNotFound) && err == nil {
		alias, err = a.gen.Generate(ctx)
		_, getErr = a.GetByAlias(ctx, alias)
	}
	return
}

func (a App) Create(ctx context.Context, url string, alias string, authorID int64, ads []int64) (links.Link, error) {
	var err error
	if alias == "" {
		alias, err = a.generateFreeAlias(ctx)
		if err != nil {
			return links.Link{}, err
		}
	}
	link := links.New(url, alias, authorID, ads)
	err = validator.Validate(link)
	if err != nil {
		return link, errors.Join(ErrInvalidContent, err)
	}
	link.ID, err = a.repo.Store(ctx, link)
	return link, err
}

func (a App) GetByID(ctx context.Context, id int64) (links.Link, error) {
	return a.repo.GetByID(ctx, id)
}

func (a App) GetByAuthor(ctx context.Context, author int64) ([]links.Link, error) {
	return a.repo.GetByAuthor(ctx, author)
}

func (a App) GetByAlias(ctx context.Context, alias string) (links.Link, error) {
	return a.repo.GetByAlias(ctx, alias)
}

func (a App) getEditable(ctx context.Context, id int64, authorID int64) (links.Link, error) {
	link, err := a.GetByID(ctx, id)
	if err != nil {
		return link, err
	}
	if link.AuthorID != authorID {
		return link, ErrPermissionDenied
	}
	return link, nil
}

func (a App) UpdateAlias(ctx context.Context, id int64, authorID int64, alias string) (links.Link, error) {
	link, err := a.getEditable(ctx, id, authorID)
	if err != nil {
		return link, err
	}
	if alias == "" {
		alias, err = a.generateFreeAlias(ctx)
		if err != nil {
			return link, err
		}
	}
	prev := link.Alias
	link.Alias = alias
	err = validator.Validate(link)
	if err != nil {
		link.Alias = prev
		return link, err
	}
	err = a.repo.UpdateAlias(ctx, id, alias)
	if err != nil {
		link.Alias = prev
	}
	return link, err
}

func (a App) AddAd(ctx context.Context, linkID int64, adID int64, authorID int64) (links.Link, error) {
	link, err := a.getEditable(ctx, linkID, authorID)
	if err != nil {
		return link, err
	}
	for i := range link.Ads {
		if link.Ads[i] == adID {
			return link, ErrAdAlreadyAdded
		}
	}
	err = a.repo.AddAd(ctx, linkID, adID)
	if err == nil {
		link.Ads = append(link.Ads, adID)
	}
	return link, err
}

func (a App) DeleteAd(ctx context.Context, linkID int64, adID int64, authorID int64) (links.Link, error) {
	link, err := a.getEditable(ctx, linkID, authorID)
	if err != nil {
		return link, err
	}
	adIdx := -1
	for i, v := range link.Ads {
		if v == adID {
			adIdx = i
			break
		}
	}
	if adIdx == -1 {
		return link, ErrNotFound
	}
	err = a.repo.DeleteAd(ctx, linkID, adID)
	if err == nil {
		link.Ads = append(link.Ads[:adIdx], link.Ads[adIdx+1:]...)
	}
	return link, err
}

func (a App) Delete(ctx context.Context, id int64, authorID int64) error {
	_, err := a.getEditable(ctx, id, authorID)
	if err != nil {
		return err
	}
	return a.repo.Delete(ctx, id)
}

func New(repo Repository, generator Generator) App {
	return App{
		repo: repo,
		gen:  generator,
	}
}

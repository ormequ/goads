package app

import (
	"context"
	"errors"
	"goads/internal/urlshortener/links"
)

type Repository interface {
	Store(ctx context.Context, link links.Link) (int64, error)
	GetByID(ctx context.Context, id int64) (links.Link, error)
	GetByAuthor(ctx context.Context, authorID int64) ([]links.Link, error)
	GetByAlias(ctx context.Context, alias string) (links.Link, error)
	UpdateAlias(ctx context.Context, id int64, alias string) error
	Delete(ctx context.Context, id int64) error
}

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

func (a App) UpdateAlias(ctx context.Context, id int64, alias string) error {
	var err error
	if alias == "" {
		alias, err = a.generateFreeAlias(ctx)
		if err != nil {
			return err
		}
	}
	return a.repo.UpdateAlias(ctx, id, alias)
}

// todo: deleteAd, addAd

func (a App) Delete(ctx context.Context, id int64) error {
	return a.repo.Delete(ctx, id)
}

func New(repo Repository, generator Generator) App {
	return App{
		repo: repo,
		gen:  generator,
	}
}

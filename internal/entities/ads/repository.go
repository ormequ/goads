package ads

import (
	"context"
	"goads/internal/entities"
)

type Repository interface {
	GetNewID(ctx context.Context) (int64, error)
	Store(ctx context.Context, ad Ad) error
	GetByID(ctx context.Context, id int64) (Ad, error)
	GetFiltered(ctx context.Context, filter entities.Filter) ([]Ad, error)
	Update(ctx context.Context, ad Ad) error
	Delete(ctx context.Context, id int64) error
}

package users

import "context"

type Repository interface {
	GetNewID(ctx context.Context) (int64, error)
	Store(ctx context.Context, user User) error
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int64) error
}

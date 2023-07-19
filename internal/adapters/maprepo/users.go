package maprepo

import (
	"context"
	"goads/internal/app"
	"goads/internal/entities/users"
	"sync"
)

type Users struct {
	byEmail sync.Map
	Repo
}

func (r *Users) Store(ctx context.Context, user users.User) (int64, error) {
	_, found := r.byEmail.Load(user.Email)
	if found {
		return -1, app.ErrAlreadyExists
	}
	id, err := r.Repo.GetNewID(ctx)
	if err != nil {
		return id, err
	}
	user.ID = id
	r.byEmail.Store(user.Email, user.ID)
	return id, r.Repo.Store(ctx, user, id)
}

func (r *Users) GetByID(ctx context.Context, id int64) (users.User, error) {
	ent, err := r.Repo.GetByID(ctx, id)
	if err != nil {
		return users.User{}, err
	}
	user, ok := ent.(users.User)
	if !ok {
		return users.User{}, app.ErrNotFound
	}
	return user, nil
}

func (r *Users) GetByEmail(ctx context.Context, email string) (users.User, error) {
	if ctx.Err() != nil {
		return users.User{}, ctx.Err()
	}
	idAny, found := r.byEmail.Load(email)
	id, ok := idAny.(int64)
	if !found || !ok {
		return users.User{}, app.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *Users) Update(ctx context.Context, user users.User) error {
	prev, err := r.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	r.byEmail.Delete(prev.Email)
	r.byEmail.Store(user.Email, user.ID)
	return r.Repo.Update(ctx, user, user.ID)
}

func NewUsers() *Users {
	return &Users{}
}

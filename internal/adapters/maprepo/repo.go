package maprepo

import (
	"context"
	"goads/internal/app"
	"sync"
	"sync/atomic"
)

type Repo struct {
	storage sync.Map
	size    atomic.Int64
}

// Store collects an entity in the map. Key - ID
func (r *Repo) Store(ctx context.Context, entity any, id int64) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	_, err := r.GetByID(ctx, id)
	if err == nil {
		return app.ErrAlreadyExists
	}
	r.storage.Store(id, entity)
	r.size.Add(1)
	return nil
}

// Update changes entity's data in map
func (r *Repo) Update(ctx context.Context, entity any, id int64) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	_, ok := r.storage.Load(id)
	if !ok {
		return app.ErrNotFound
	}
	r.storage.Store(id, entity)
	return nil
}

// GetNewID returns the ID for new entity
func (r *Repo) GetNewID(ctx context.Context) (int64, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	return r.size.Load(), nil
}

// GetByID returns an entity found by ID or the error if not found
func (r *Repo) GetByID(ctx context.Context, id int64) (any, error) {
	v, okLoad := r.storage.Load(id)
	if !okLoad {
		return v, app.ErrNotFound
	}
	if e := ctx.Err(); e != nil {
		return v, e
	}
	return v, nil
}

// Delete removes entity by got id from storage
func (r *Repo) Delete(ctx context.Context, id int64) error {
	if _, ok := r.storage.Load(id); !ok {
		return app.ErrNotFound
	}
	r.storage.Delete(id)
	if e := ctx.Err(); e != nil {
		return e
	}
	return nil
}

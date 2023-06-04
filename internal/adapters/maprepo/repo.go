package maprepo

import (
	"context"
	"goads/internal/app"
	"goads/internal/entities"
	"sync"
	"sync/atomic"
)

type Repo[T entities.Interface] struct {
	storage sync.Map
	size    atomic.Int64
}

// Store collects an entity in the map. Key - ID
func (m *Repo[T]) Store(ctx context.Context, entity T) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	_, err := m.GetByID(ctx, entity.GetID())
	if err == nil {
		return app.ErrAlreadyExists
	}
	m.storage.Store(entity.GetID(), entity)
	m.size.Add(1)
	return nil
}

// GetFiltered returns list of ads where for each entity applied functions in filter returns true
func (m *Repo[T]) GetFiltered(ctx context.Context, filter entities.Filter) ([]T, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	res := make([]T, 0, m.size.Load())
	m.storage.Range(func(key, value any) bool {
		val := value.(T)
		for _, f := range filter {
			if !f(val) {
				return true
			}
		}
		res = append(res, val)
		return true
	})
	return res, nil
}

// Update changes entity's data in map
func (m *Repo[T]) Update(ctx context.Context, entity T) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	_, ok := m.storage.Load(entity.GetID())
	if !ok {
		return app.ErrNotFound
	}
	m.storage.Store(entity.GetID(), entity)
	return nil
}

// GetNewID returns the ID for new entity
func (m *Repo[T]) GetNewID(ctx context.Context) (int64, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	return m.size.Load(), nil
}

// GetByID returns an entity found by ID or the error if not found
func (m *Repo[T]) GetByID(ctx context.Context, id int64) (T, error) {
	v, okLoad := m.storage.Load(id)
	t, okConv := v.(T)
	if !okLoad || !okConv {
		return t, app.ErrNotFound
	}
	if e := ctx.Err(); e != nil {
		return t, e
	}
	return t, nil
}

// Delete removes entity by got id from storage
func (m *Repo[T]) Delete(ctx context.Context, id int64) error {
	if _, ok := m.storage.Load(id); !ok {
		return app.ErrNotFound
	}
	m.storage.Delete(id)
	if e := ctx.Err(); e != nil {
		return e
	}
	return nil
}

func New[T entities.Interface]() *Repo[T] {
	return &Repo[T]{}
}

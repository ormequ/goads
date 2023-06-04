package maprepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"goads/internal/app"
	"goads/internal/entities"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"testing"
)

func TestRepo_StoreGet(t *testing.T) {
	type testCase[T entities.Interface] struct {
		name    string
		ctx     context.Context
		entity  T
		wantErr error
	}

	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := [...]testCase[ads.Ad]{
		{
			name:   "Valid ad",
			ctx:    ctx,
			entity: ads.New(0, "hello", "world", 0),
		},
		{
			name:   "Valid ad",
			ctx:    ctx,
			entity: ads.New(1, "test", "2", 0),
		},
		{
			name:   "Valid ad",
			ctx:    ctx,
			entity: ads.New(2, "test", "3", 0),
		},
		{
			name:    "Already exists",
			ctx:     ctx,
			entity:  ads.New(0, "hello", "world", 0),
			wantErr: app.ErrAlreadyExists,
		},
		{
			name:    "Canceled ctx",
			ctx:     canceledCtx,
			wantErr: canceledCtx.Err(),
		},
	}

	repo := New[ads.Ad]()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Store(tt.ctx, tt.entity)
			assert.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				got, err := repo.GetByID(tt.ctx, tt.entity.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.entity, got)
			}
		})
	}
	_, err := repo.GetByID(canceledCtx, 0)
	assert.ErrorIs(t, err, canceledCtx.Err())
}

func TestRepo_GetNewID(t *testing.T) {
	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	usersList := [...]users.User{
		users.New(0, "1", "1"),
		users.New(1, "2", "2"),
		users.New(2, "3", "3"),
		users.New(3, "4", "4"),
		users.New(4, "5", "5"),
	}
	repo := New[users.User]()

	for _, user := range usersList {
		id, err := repo.GetNewID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, id)
		err = repo.Store(ctx, user)
		assert.NoError(t, err)
	}
	_, err := repo.GetNewID(canceledCtx)
	assert.ErrorIs(t, err, canceledCtx.Err())
}

func TestRepo_Delete(t *testing.T) {
	type testCase[T entities.Interface] struct {
		name    string
		ctx     context.Context
		id      int64
		wantErr error
	}
	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	adsList := [...]ads.Ad{
		ads.New(0, "1", "1", 0),
		ads.New(1, "2", "2", 0),
		ads.New(2, "3", "3", 0),
		ads.New(3, "4", "4", 0),
		ads.New(4, "5", "5", 0),
	}
	repo := New[ads.Ad]()
	for _, ad := range adsList {
		assert.NoError(t, repo.Store(ctx, ad))
	}

	tests := [...]testCase[ads.Ad]{
		{
			name: "Valid",
			ctx:  ctx,
			id:   adsList[0].ID,
		},
		{
			name: "Valid",
			ctx:  ctx,
			id:   adsList[4].ID,
		},
		{
			name:    "Canceled context",
			ctx:     canceledCtx,
			id:      adsList[1].ID,
			wantErr: canceledCtx.Err(),
		},
		{
			name:    "Already deleted",
			ctx:     ctx,
			id:      adsList[0].ID,
			wantErr: app.ErrNotFound,
		},
		{
			name:    "Not found",
			ctx:     ctx,
			id:      100,
			wantErr: app.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.ctx, tt.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRepo_Update(t *testing.T) {
	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	adsList := [...]ads.Ad{
		ads.New(0, "1", "1", 0),
		ads.New(1, "2", "2", 0),
		ads.New(2, "3", "3", 0),
	}
	repo := New[ads.Ad]()
	for _, ad := range adsList {
		assert.NoError(t, repo.Store(ctx, ad))
		ad.Text = "hello world"
		assert.NoError(t, repo.Update(ctx, ad))
	}
	assert.ErrorIs(t, repo.Update(ctx, ads.New(123, "_", "_", 0)), app.ErrNotFound)
	assert.ErrorIs(t, repo.Update(canceledCtx, adsList[0]), canceledCtx.Err())
}

func TestRepo_GetFiltered(t *testing.T) {
	type testCase[T entities.Interface] struct {
		name    string
		ctx     context.Context
		filter  entities.Filter
		len     int
		res     ads.Ad
		wantErr error
	}
	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	adsList := [...]ads.Ad{
		ads.New(0, "1", "1", 0),
		ads.New(1, "2", "2", 0),
		ads.New(2, "3", "3", 0),
		ads.New(3, "4", "4", 0),
		ads.New(4, "5", "5", 0),
	}
	repo := New[ads.Ad]()
	for _, ad := range adsList {
		assert.NoError(t, repo.Store(ctx, ad))
	}

	tests := [...]testCase[ads.Ad]{
		{
			name: "Correct filter",
			ctx:  ctx,
			filter: entities.Filter{
				func(i entities.Interface) bool {
					return i.GetID() > 0
				},
			},
			len: 4,
		},
		{
			name:    "Canceled context",
			ctx:     canceledCtx,
			filter:  entities.Filter{},
			len:     0,
			wantErr: canceledCtx.Err(),
		},
		{
			name: "Correct filter",
			ctx:  ctx,
			filter: entities.Filter{
				func(i entities.Interface) bool {
					return i.GetID() == 3
				},
			},
			len: 1,
			res: adsList[3],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetFiltered(tt.ctx, tt.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Len(t, got, tt.len)
			if len(got) == 1 {
				assert.Equal(t, tt.res, got[0])
			}
		})
	}
}

func BenchmarkRepo_Store(b *testing.B) {
	repo := New[ads.Ad]()
	for i := 0; i < b.N; i++ {
		ad := ads.New(int64(i), "hello", "world", 0)
		err := repo.Store(context.Background(), ad)
		assert.NoError(b, err)
	}
}

func BenchmarkRepo_Get(b *testing.B) {
	repo := New[ads.Ad]()
	for i := 0; i < b.N; i++ {
		ad := ads.New(int64(i), "hello", "world", 0)
		err := repo.Store(context.Background(), ad)
		assert.NoError(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetByID(context.Background(), int64(i))
		assert.NoError(b, err)
	}
}

func BenchmarkRepo_Delete(b *testing.B) {
	repo := New[ads.Ad]()
	for i := 0; i < b.N; i++ {
		ad := ads.New(int64(i), "hello", "world", 0)
		err := repo.Store(context.Background(), ad)
		assert.NoError(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.Delete(context.Background(), int64(i))
		assert.NoError(b, err)
	}
}

func FuzzRepo_GetDelete(f *testing.F) {
	repo := New[ads.Ad]()
	ctx := context.Background()
	f.Fuzz(func(t *testing.T, id int64, title string, text string, authorID int64) {
		expected := ads.New(id, title, text, authorID)
		err := repo.Store(ctx, expected)
		assert.NoError(t, err)
		got, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		err = repo.Delete(ctx, id)
		assert.NoError(t, err)
	})
}

package ad_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goads/internal/app"
	"goads/internal/app/ad"
	"goads/internal/app/ad/mocks"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"testing"
)

func TestAds_CreateDelete(t *testing.T) {
	type args struct {
		ctx      context.Context
		title    string
		text     string
		authorID int64
	}

	getter := mocks.NewUsersGetter(t)
	getter.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) (users.User, error) {
			if id == 321 {
				return users.User{}, app.ErrNotFound
			}
			return users.User{}, nil
		})

	storeRepo := mocks.NewRepository(t)
	storeRepo.
		On("Store", mock.Anything, mock.AnythingOfType("ads.Ad")).
		Return(func(_ context.Context, ad ads.Ad) (int64, error) {
			if ad.AuthorID == 123 {
				return 0, app.ErrAlreadyExists
			}
			return 0, nil
		})
	delRepo := mocks.NewRepository(t)
	delRepo.
		On("Store", mock.Anything, mock.AnythingOfType("ads.Ad")).
		Return(int64(0), nil)
	delRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(ads.Ad{AuthorID: 1}, nil)
	delRepo.
		On("Delete", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(ctx context.Context, id int64) error {
			if ctx != nil {
				return app.ErrNotFound
			}
			return nil
		})

	tests := [...]struct {
		name      string
		args      args
		want      ads.Ad
		storeErr  error
		deleteErr error
		repo      *mocks.Repository
	}{
		{
			name: "Standard ad",
			args: args{title: "hello", text: "world", authorID: 1},
			want: ads.Ad{
				AuthorID:  1,
				Published: false,
				Title:     "hello",
				Text:      "world",
			},
			repo: delRepo,
		},
		{
			name: "Delete not found",
			args: args{ctx: context.Background(), title: "hello", text: "world", authorID: 1},
			want: ads.Ad{
				AuthorID:  1,
				Published: false,
				Title:     "hello",
				Text:      "world",
			},
			deleteErr: app.ErrNotFound,
			repo:      delRepo,
		},
		{
			name:     "Invalid content",
			args:     args{title: "", text: "world", authorID: 1},
			storeErr: app.ErrInvalidContent,
			repo:     storeRepo,
		},
		{
			name:     "Already exists",
			args:     args{title: "hello", text: "world", authorID: 123},
			storeErr: app.ErrAlreadyExists,
			repo:     storeRepo,
		},
		{
			name:     "User doesn't exist",
			args:     args{title: "hello", text: "world", authorID: 321},
			storeErr: app.ErrNotFound,
			repo:     storeRepo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ad.New(tt.repo, getter)
			got, err := a.Create(tt.args.ctx, tt.args.title, tt.args.text, tt.args.authorID)
			assert.ErrorIs(t, err, tt.storeErr)
			if err != nil {
				return
			}
			assert.Zero(t, got.ID)
			assert.Equal(t, tt.want.AuthorID, got.AuthorID)
			assert.Equal(t, tt.want.Published, got.Published)
			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Text, got.Text)

			err = a.Delete(tt.args.ctx, got.ID, got.AuthorID)
			assert.ErrorIs(t, err, tt.deleteErr)

		})
	}
}

func TestAds_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}

	repo := mocks.NewRepository(t)
	repo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) (ads.Ad, error) {
			if id == 123 {
				return ads.Ad{}, app.ErrNotFound
			}
			return ads.Ad{}, nil
		})

	tests := [...]struct {
		name string
		args args
		err  error
	}{
		{
			name: "Valid getting",
			args: args{},
		},
		{
			name: "Not found",
			args: args{id: 123},
			err:  app.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ad.New(repo, nil)
			_, err := a.GetByID(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestAds_ChangeStatus(t *testing.T) {
	type args struct {
		id        int64
		userID    int64
		published bool
	}

	getter := mocks.NewUsersGetter(t)
	getter.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{}, nil)

	getByID := func(_ context.Context, id int64) (ads.Ad, error) {
		if id == 123 {
			return ads.Ad{}, app.ErrNotFound
		}
		return ads.Ad{
			Title:    "hello",
			Text:     "world",
			AuthorID: id,
		}, nil
	}

	repo := mocks.NewRepository(t)
	repo.
		On("Update", mock.Anything, mock.AnythingOfType("ads.Ad")).
		Return(nil)
	repo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(getByID)

	errRepo := mocks.NewRepository(t)
	errRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(getByID)

	tests := [...]struct {
		name   string
		args   args
		err    error
		repo   *mocks.Repository
		getter *mocks.UsersGetter
	}{
		{
			name:   "Valid changing",
			args:   args{},
			repo:   repo,
			getter: getter,
		},
		{
			name: "Ad not found",
			args: args{id: 123},
			err:  app.ErrNotFound,
			repo: repo,
		},
		{
			name:   "Permission denied",
			args:   args{userID: 123},
			err:    app.ErrPermissionDenied,
			repo:   errRepo,
			getter: getter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ad.New(tt.repo, tt.getter)
			assert.ErrorIs(t, a.ChangeStatus(context.Background(), tt.args.id, tt.args.userID, tt.args.published), tt.err)
		})
	}
}

func TestAds_GetFiltered(t *testing.T) {
	type args struct {
		ctx context.Context
		opt ad.Filter
	}

	getter := mocks.NewUsersGetter(t)
	getter.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) (users.User, error) {
			if id == 321 {
				return users.User{}, app.ErrNotFound
			}
			return users.User{}, nil
		})

	repo := mocks.NewRepository(t)
	repo.
		On("GetFiltered", mock.Anything, mock.AnythingOfType("ad.Filter")).
		Return(func(ctx context.Context, _ ad.Filter) ([]ads.Ad, error) {
			if ctx != nil {
				return nil, app.ErrInvalidFilter
			}
			return []ads.Ad{}, nil
		})

	tests := [...]struct {
		name string
		args args
		err  error
		repo *mocks.Repository
	}{
		{
			name: "Valid filter",
			args: args{},
			repo: repo,
		},
		{
			name: "User not found",
			args: args{opt: ad.Filter{AuthorID: 321}},
			err:  app.ErrInvalidFilter,
			repo: mocks.NewRepository(t),
		},
		{
			name: "Repository error",
			args: args{ctx: context.Background()},
			err:  app.ErrInvalidFilter,
			repo: repo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ad.New(tt.repo, getter)
			_, err := a.GetFiltered(tt.args.ctx, tt.args.opt)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestAds_Search(t *testing.T) {
	type args struct {
		ctx   context.Context
		title string
	}

	repo := mocks.NewRepository(t)
	repo.
		On("GetFiltered", mock.Anything, mock.AnythingOfType("ad.Filter")).
		Return(func(ctx context.Context, _ ad.Filter) ([]ads.Ad, error) {
			if ctx != nil {
				return nil, app.ErrInvalidFilter
			}
			return []ads.Ad{}, nil
		})

	tests := [...]struct {
		name string
		args args
		err  error
		repo *mocks.Repository
	}{
		{
			name: "Valid filter",
			args: args{},
			repo: repo,
		},
		{
			name: "Repository error",
			args: args{ctx: context.Background()},
			err:  app.ErrInvalidFilter,
			repo: repo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ad.New(tt.repo, nil)
			_, err := a.Search(tt.args.ctx, tt.args.title)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestAds_Update(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     int64
		userID int64
		title  string
		text   string
	}

	getter := mocks.NewUsersGetter(t)
	getter.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) (users.User, error) {
			if id == 321 {
				return users.User{}, app.ErrNotFound
			}
			return users.User{}, nil
		})

	getByID := func(_ context.Context, id int64) (ads.Ad, error) {
		if id == 123 {
			return ads.Ad{}, app.ErrNotFound
		}
		return ads.Ad{
			Title:    "hello",
			Text:     "world",
			AuthorID: 1,
		}, nil
	}

	repo := mocks.NewRepository(t)
	repo.
		On("Update", mock.Anything, mock.AnythingOfType("ads.Ad")).
		Return(func(_ context.Context, ad ads.Ad) error {
			if ad.AuthorID == 123 {
				return app.ErrNotFound
			}
			return nil
		})
	repo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(getByID)

	errRepo := mocks.NewRepository(t)
	errRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(getByID)

	tests := [...]struct {
		name   string
		args   args
		err    error
		repo   *mocks.Repository
		getter *mocks.UsersGetter
	}{
		{
			name:   "Valid changing",
			args:   args{userID: 1, title: "hello", text: "world"},
			repo:   repo,
			getter: getter,
		},
		{
			name: "Ad not found",
			args: args{id: 123, userID: 1, title: "hello", text: "world"},
			err:  app.ErrNotFound,
			repo: repo,
		},
		{
			name:   "Permission denied",
			args:   args{userID: 123, title: "hello", text: "world"},
			err:    app.ErrPermissionDenied,
			repo:   errRepo,
			getter: getter,
		},
		{
			name:   "Invalid update",
			args:   args{userID: 1, title: "", text: "world"},
			err:    app.ErrInvalidContent,
			repo:   errRepo,
			getter: getter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ad.New(tt.repo, tt.getter)
			assert.ErrorIs(t, a.Update(tt.args.ctx, tt.args.id, tt.args.userID, tt.args.title, tt.args.text), tt.err)
		})
	}
}

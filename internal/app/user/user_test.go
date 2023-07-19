package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"goads/internal/app"
	"goads/internal/app/user/mocks"
	"goads/internal/entities/users"
	"testing"
)

func TestUsers_CreateDelete(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx   context.Context
		email string
		name  string
	}

	storeRepo := mocks.NewRepository(t)
	storeRepo.
		On("Store", mock.Anything, mock.AnythingOfType("users.User")).
		Return(func(_ context.Context, u users.User) (int64, error) {
			if u.Name == "error" {
				return 0, app.ErrAlreadyExists
			}
			return 0, nil
		})

	delRepo := mocks.NewRepository(t)
	delRepo.
		On("Store", mock.Anything, mock.AnythingOfType("users.User")).
		Return(int64(0), nil)
	delRepo.
		On("Delete", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(ctx context.Context, _ int64) error {
			if ctx != nil {
				return app.ErrNotFound
			}
			return nil
		})

	tests := [...]struct {
		name     string
		fields   fields
		args     args
		want     users.User
		storeErr error
		delErr   error
	}{
		{
			name:   "Valid user",
			fields: fields{delRepo},
			args:   args{email: "a@d.c", name: "test"},
			want: users.User{
				ID:    0,
				Email: "a@d.c",
				Name:  "test",
			},
		},
		{
			name:     "Invalid user data",
			fields:   fields{storeRepo},
			args:     args{email: "a@d.c"},
			storeErr: app.ErrInvalidContent,
		},
		{
			name:     "Store error",
			fields:   fields{storeRepo},
			args:     args{name: "error", email: "error@test.com"},
			storeErr: app.ErrAlreadyExists,
		},
		{
			name:   "Delete error",
			fields: fields{delRepo},
			args:   args{ctx: context.Background(), name: "test", email: "test@test.com"},
			want: users.User{
				ID:    0,
				Name:  "test",
				Email: "test@test.com",
			},
			delErr: app.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := App{
				repository: tt.fields.repository,
			}
			got, err := u.Create(tt.args.ctx, tt.args.email, tt.args.name)
			assert.ErrorIs(t, err, tt.storeErr)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
			err = u.Delete(tt.args.ctx, got.ID)
			assert.ErrorIs(t, err, tt.delErr)
		})
	}
}

func TestUsers_ChangeEmail(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx   context.Context
		id    int64
		email string
	}

	getRepo := mocks.NewRepository(t)
	getRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{Name: "test"}, func(_ context.Context, id int64) error {
			if id == 123 {
				return app.ErrNotFound
			}
			return nil
		})

	updateRepo := mocks.NewRepository(t)
	updateRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{Name: "test"}, nil)
	updateRepo.
		On("Update", mock.Anything, mock.AnythingOfType("users.User")).
		Return(func(_ context.Context, u users.User) error {
			if u.Email == "error@test.com" {
				return app.ErrNotFound
			}
			return nil
		})

	tests := [...]struct {
		name   string
		fields fields
		args   args
		err    error
	}{
		{
			name:   "Valid changing",
			fields: fields{updateRepo},
			args:   args{email: "test@test.com"},
		},
		{
			name:   "Invalid email",
			fields: fields{getRepo},
			err:    app.ErrInvalidContent,
		},
		{
			name:   "Get error",
			fields: fields{getRepo},
			args:   args{id: 123},
			err:    app.ErrNotFound,
		},
		{
			name:   "Update error",
			fields: fields{updateRepo},
			args:   args{email: "error@test.com"},
			err:    app.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := App{
				repository: tt.fields.repository,
			}
			err := u.ChangeEmail(tt.args.ctx, tt.args.id, tt.args.email)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestUsers_ChangeName(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx  context.Context
		id   int64
		name string
	}

	getRepo := mocks.NewRepository(t)
	getRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{Email: "test@test.com"}, func(_ context.Context, id int64) error {
			if id == 123 {
				return app.ErrNotFound
			}
			return nil
		})

	updateRepo := mocks.NewRepository(t)
	updateRepo.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{Email: "test@test.com"}, nil)
	updateRepo.
		On("Update", mock.Anything, mock.AnythingOfType("users.User")).
		Return(func(_ context.Context, u users.User) error {
			if u.Name == "error" {
				return app.ErrNotFound
			}
			return nil
		})

	tests := [...]struct {
		name   string
		fields fields
		args   args
		err    error
	}{
		{
			name:   "Valid changing",
			fields: fields{updateRepo},
			args:   args{name: "test"},
		},
		{
			name:   "Invalid name",
			fields: fields{getRepo},
			err:    app.ErrInvalidContent,
		},
		{
			name:   "Get error",
			fields: fields{getRepo},
			args:   args{id: 123},
			err:    app.ErrNotFound,
		},
		{
			name:   "Update error",
			fields: fields{updateRepo},
			args:   args{name: "error"},
			err:    app.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := App{
				repository: tt.fields.repository,
			}
			err := u.ChangeName(tt.args.ctx, tt.args.id, tt.args.name)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

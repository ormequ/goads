package app

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goads/internal/auth/app/mocks"
	"goads/internal/auth/users"
	"testing"
)

func storeRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("Store", mock.Anything, mock.AnythingOfType("users.User")).
		Return(func(ctx context.Context, user users.User) (int64, error) {
			if ctx.Err() != nil {
				return -1, ctx.Err()
			}
			if user.Email == "already@exists.com" {
				return -1, ErrEmailAlreadyExists
			}
			return 0, nil
		})
	return r
}

func getByEmailRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
		Return(func(ctx context.Context, email string) (users.User, error) {
			if ctx.Err() != nil {
				return users.User{}, ctx.Err()
			}
			if email == "incorrect@credentials.com" {
				return users.User{}, ErrIncorrectCredentials
			}
			return users.User{Email: email}, nil
		})
	return r
}

func getByIDUpdateRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(users.User{ID: 0, Name: "test", Email: "test@test.com"}, nil)
	r.
		On("Update", mock.Anything, mock.AnythingOfType("users.User")).
		Return(func(ctx context.Context, user users.User) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if user.Password == "not_found_hash" {
				return ErrNotFound
			}
			return nil
		})

	return r
}

func getByIDRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(ctx context.Context, id int64) (users.User, error) {
			if ctx.Err() != nil {
				return users.User{}, ctx.Err()
			}
			if id == -1 {
				return users.User{}, ErrNotFound
			}
			return users.User{ID: id, Name: "test", Email: "test@test.com"}, nil
		})

	return r
}

func hashGenerator(t *testing.T) Hasher {
	h := mocks.NewHasher(t)
	h.
		On("Generate", mock.Anything, mock.AnythingOfType("string")).
		Return(func(ctx context.Context, pass string) (string, error) {
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			if pass == "" {
				return "", ErrPasswordToShort
			}
			return fmt.Sprintf("%s_hash", pass), nil
		})
	return h
}

func hashComparator(t *testing.T) Hasher {
	h := mocks.NewHasher(t)
	h.
		On("Compare", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(func(ctx context.Context, hash string, pass string) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if pass == "" {
				return ErrIncorrectCredentials
			}
			return nil
		})
	return h
}

func tokenizer(t *testing.T) Tokenizer {
	tok := mocks.NewTokenizer(t)
	tok.
		On("Generate", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(ctx context.Context, id int64) (string, error) {
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return "token", nil
		})
	return tok

}

func TestApp_Register(t *testing.T) {
	type fields struct {
		Repo      Repository
		Tokenizer Tokenizer
		Hasher    Hasher
		Validator Validator
	}
	type args struct {
		ctx      context.Context
		email    string
		name     string
		password string
	}

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    users.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct registering",
			fields: fields{
				Repo:   storeRepo(t),
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@test.com",
				name:     "test",
				password: "asdf",
			},
			want: users.User{
				ID:       0,
				Email:    "test@test.com",
				Name:     "test",
				Password: "asdf_hash",
			},
		},
		{
			name: "validate error",
			fields: fields{
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx:      context.Background(),
				email:    "",
				name:     "test",
				password: "asdf",
			},
			want: users.User{
				ID:       0,
				Email:    "",
				Name:     "test",
				Password: "asdf_hash",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrInvalidContent, i)
			},
		},
		{
			name: "repository error",
			fields: fields{
				Repo:   storeRepo(t),
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx:      context.Background(),
				email:    "already@exists.com",
				name:     "test",
				password: "asdf",
			},
			want: users.User{
				ID:       -1,
				Email:    "already@exists.com",
				Name:     "test",
				Password: "asdf_hash",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmailAlreadyExists, i)
			},
		},
		{
			name: "canceled ctx",
			fields: fields{
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx: canceledCtx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, canceledCtx.Err(), i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo:      tt.fields.Repo,
				Tokenizer: tt.fields.Tokenizer,
				Hasher:    tt.fields.Hasher,
				Validator: tt.fields.Validator,
			}
			got, err := a.Register(tt.args.ctx, tt.args.email, tt.args.name, tt.args.password)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("Register(%v, %v, %v, %v)", tt.args.ctx, tt.args.email, tt.args.name, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Register(%v, %v, %v, %v)", tt.args.ctx, tt.args.email, tt.args.name, tt.args.password)
		})
	}
}

func TestApp_Authenticate(t *testing.T) {
	type fields struct {
		Repo      Repository
		Tokenizer Tokenizer
		Hasher    Hasher
		Validator Validator
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct auth",
			fields: fields{
				Repo:      getByEmailRepo(t),
				Tokenizer: tokenizer(t),
				Hasher:    hashComparator(t),
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@test.com",
				password: "test",
			},
			want: "token",
		},
		{
			name: "canceled ctx",
			fields: fields{
				Repo: getByEmailRepo(t),
			},
			args: args{
				ctx: canceledCtx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, canceledCtx.Err(), i)
			},
		},
		{
			name: "incorrect email",
			fields: fields{
				Repo: getByEmailRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				email:    "incorrect@credentials.com",
				password: "test",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrIncorrectCredentials, i)
			},
		},
		{
			name: "incorrect password",
			fields: fields{
				Repo:   getByEmailRepo(t),
				Hasher: hashComparator(t),
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@test.com",
				password: "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrIncorrectCredentials, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo:      tt.fields.Repo,
				Tokenizer: tt.fields.Tokenizer,
				Hasher:    tt.fields.Hasher,
				Validator: tt.fields.Validator,
			}
			got, err := a.Authenticate(tt.args.ctx, tt.args.email, tt.args.password)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("Authenticate(%v, %v, %v)", tt.args.ctx, tt.args.email, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Authenticate(%v, %v, %v)", tt.args.ctx, tt.args.email, tt.args.password)
		})
	}
}

func TestApp_ChangePassword(t *testing.T) {
	type fields struct {
		Repo      Repository
		Tokenizer Tokenizer
		Hasher    Hasher
		Validator Validator
	}
	type args struct {
		ctx      context.Context
		id       int64
		password string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    users.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct changing",
			fields: fields{
				Repo:   getByIDUpdateRepo(t),
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx:      context.Background(),
				id:       0,
				password: "test",
			},
			want: users.User{
				ID:       0,
				Name:     "test",
				Email:    "test@test.com",
				Password: "test_hash",
			},
		},
		{
			name: "not found by id",
			fields: fields{
				Repo:   getByIDRepo(t),
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx:      context.Background(),
				id:       -1,
				password: "test",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
		{
			name: "new password to short",
			fields: fields{
				Hasher: hashGenerator(t),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrPasswordToShort, i)
			},
		},
		{
			name: "not found on update",
			fields: fields{
				Hasher: hashGenerator(t),
				Repo:   getByIDUpdateRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				id:       0,
				password: "not_found",
			},
			want: users.User{
				ID:       0,
				Email:    "test@test.com",
				Name:     "test",
				Password: "not_found_hash",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo:      tt.fields.Repo,
				Tokenizer: tt.fields.Tokenizer,
				Hasher:    tt.fields.Hasher,
				Validator: tt.fields.Validator,
			}
			got, err := a.ChangePassword(tt.args.ctx, tt.args.id, tt.args.password)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("ChangePassword(%v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ChangePassword(%v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.password)
		})
	}
}

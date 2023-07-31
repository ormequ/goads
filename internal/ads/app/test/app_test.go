package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goads/internal/ads/ads"
	"goads/internal/ads/app"
	"goads/internal/ads/app/mocks"
	"testing"
	"time"
)

func getByIDRepo(t *testing.T) app.Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) (ads.Ad, error) {
			if id == -1 {
				return ads.Ad{}, app.ErrAdNotFound
			}
			return ads.Ad{
				Title: "test title",
				Text:  "test text",
			}, nil
		})
	return r
}

func getByIDUpdateRepo(t *testing.T) app.Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(ads.Ad{
			Title: "test title",
			Text:  "test text",
		}, nil)

	r.
		On("Update", mock.Anything, mock.AnythingOfType("ads.Ad")).
		Return(func(_ context.Context, ad ads.Ad) error {
			if ad.Title == "not found" {
				return app.ErrAdNotFound
			}
			return nil
		})

	return r
}

func getByIDDeleteRepo(t *testing.T) app.Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(ads.Ad{
			Title: "test title",
			Text:  "test text",
		}, nil)

	r.
		On("Delete", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, id int64) error {
			if id == -1 {
				return app.ErrAdNotFound
			}
			return nil
		})

	return r
}

func storeRepo(t *testing.T) app.Repository {
	r := mocks.NewRepository(t)
	r.
		On("Store", mock.Anything, mock.AnythingOfType("ads.Ad")).
		Return(func(_ context.Context, ad ads.Ad) (int64, error) {
			if ad.Title == "no author" {
				return -1, app.ErrAuthorNotFound
			}
			return 0, nil
		})
	return r
}

func TestApp_Create(t *testing.T) {
	type fields struct {
		repository app.Repository
	}
	type args struct {
		ctx      context.Context
		title    string
		text     string
		authorID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ads.Ad
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "correct creating",
			fields: fields{repository: storeRepo(t)},
			args: args{
				ctx:   context.Background(),
				title: "correct",
				text:  "correct",
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  true,
				Title:      "correct",
				Text:       "correct",
				CreateDate: time.Now().UTC(),
				UpdateDate: time.Now().UTC(),
			},
		},
		{
			name:   "author not found",
			fields: fields{repository: storeRepo(t)},
			args: args{
				ctx:   context.Background(),
				title: "no author",
				text:  "valid",
			},
			want: ads.Ad{
				ID:         -1,
				AuthorID:   0,
				Published:  true,
				Title:      "no author",
				Text:       "valid",
				CreateDate: time.Now().UTC(),
				UpdateDate: time.Now().UTC(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrAuthorNotFound)
			},
		},
		{
			name:   "invalid title",
			fields: fields{},
			args: args{
				ctx:   context.Background(),
				title: "",
				text:  "valid",
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  true,
				Title:      "",
				Text:       "valid",
				CreateDate: time.Now().UTC(),
				UpdateDate: time.Now().UTC(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrInvalidContent)
			},
		},
		{
			name:   "invalid text",
			fields: fields{},
			args: args{
				ctx:   context.Background(),
				title: "valid",
				text:  "",
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  true,
				Title:      "valid",
				Text:       "",
				CreateDate: time.Now().UTC(),
				UpdateDate: time.Now().UTC(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrInvalidContent, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := app.App{
				Repo: tt.fields.repository,
			}
			got, err := a.Create(tt.args.ctx, tt.args.title, tt.args.text, tt.args.authorID)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("Create(%v, %v, %v, %v)", tt.args.ctx, tt.args.title, tt.args.text, tt.args.authorID)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			tt.want.UpdateDate = tt.want.UpdateDate.Truncate(time.Second)
			tt.want.CreateDate = tt.want.CreateDate.Truncate(time.Second)
			got.UpdateDate = got.UpdateDate.Truncate(time.Second)
			got.CreateDate = got.CreateDate.Truncate(time.Second)
			assert.Equalf(t, tt.want, got, "Create(%v, %v, %v, %v)", tt.args.ctx, tt.args.title, tt.args.text, tt.args.authorID)
		})
	}
}

func TestApp_Delete(t *testing.T) {
	type fields struct {
		repository app.Repository
	}
	type args struct {
		ctx    context.Context
		id     int64
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "correct deleting",
			fields: fields{repository: getByIDDeleteRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 0,
			},
		},
		{
			name:   "not found",
			fields: fields{repository: getByIDRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     -1,
				userID: 0,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrAdNotFound, i)
			},
		},
		{
			name:   "permission denied",
			fields: fields{repository: getByIDRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 1,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrPermissionDenied, i)
			},
		},
		{
			name:   "deleting error",
			fields: fields{repository: getByIDDeleteRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     -1,
				userID: 0,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrAdNotFound, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := app.App{
				Repo: tt.fields.repository,
			}
			err := a.Delete(tt.args.ctx, tt.args.id, tt.args.userID)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				tt.wantErr(t, err, fmt.Sprintf("Delete(%v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.userID))
			}
		})
	}
}

func TestApp_Update(t *testing.T) {
	type fields struct {
		repository app.Repository
	}
	type args struct {
		ctx    context.Context
		id     int64
		userID int64
		title  string
		text   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ads.Ad
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "correct updating",
			fields: fields{repository: getByIDUpdateRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 0,
				title:  "valid",
				text:   "valid",
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  false,
				Title:      "valid",
				Text:       "valid",
				CreateDate: time.Time{},
				UpdateDate: time.Now().UTC(),
			},
		},
		{
			name:   "not found",
			fields: fields{repository: getByIDRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     -1,
				userID: 0,
				title:  "valid",
				text:   "valid",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrAdNotFound, i)
			},
			want: ads.Ad{},
		},
		{
			name:   "permission denied",
			fields: fields{repository: getByIDRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 1,
				title:  "valid",
				text:   "valid",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrPermissionDenied, i)
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  false,
				Title:      "test title",
				Text:       "test text",
				CreateDate: time.Time{},
				UpdateDate: time.Time{},
			},
		},
		{
			name:   "invalid title",
			fields: fields{repository: getByIDRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 0,
				title:  "",
				text:   "valid",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrInvalidContent, i)
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  false,
				Title:      "test title",
				Text:       "test text",
				CreateDate: time.Time{},
				UpdateDate: time.Time{},
			},
		},
		{
			name:   "invalid text",
			fields: fields{repository: getByIDRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 0,
				title:  "valid",
				text:   "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrInvalidContent, i)
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  false,
				Title:      "test title",
				Text:       "test text",
				CreateDate: time.Time{},
				UpdateDate: time.Time{},
			},
		},
		{
			name:   "updating error",
			fields: fields{repository: getByIDUpdateRepo(t)},
			args: args{
				ctx:    context.Background(),
				id:     0,
				userID: 0,
				title:  "not found",
				text:   "valid",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, app.ErrAdNotFound, i)
			},
			want: ads.Ad{
				ID:         0,
				AuthorID:   0,
				Published:  false,
				Title:      "test title",
				Text:       "test text",
				CreateDate: time.Time{},
				UpdateDate: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := app.App{
				Repo: tt.fields.repository,
			}
			got, err := a.Update(tt.args.ctx, tt.args.id, tt.args.userID, tt.args.title, tt.args.text)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("Update(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.userID, tt.args.title, tt.args.text)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "Update(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.userID, tt.args.title, tt.args.text)
		})
	}
}

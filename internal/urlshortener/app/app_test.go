package app

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goads/internal/urlshortener/app/mocks"
	"goads/internal/urlshortener/entities/ads"
	"goads/internal/urlshortener/entities/links"
	"goads/internal/urlshortener/entities/redirects"
	"testing"
)

func storeRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("Store", mock.Anything, mock.AnythingOfType("links.Link")).
		Return(func(_ context.Context, link links.Link) (int64, error) {
			if link.URL == "https://already-exists.com" {
				return 0, ErrAlreadyExists
			}
			return 0, nil
		})
	return r
}

func adsService(t *testing.T) AdsService {
	a := mocks.NewAdsService(t)
	a.
		On("GetOnlyPublished", mock.Anything, mock.AnythingOfType("[]int64")).
		Return(func(_ context.Context, ids []int64) ([]ads.Ad, error) {
			if len(ids) == 3 {
				ids = []int64{ids[0], ids[2]}
			}
			res := make([]ads.Ad, len(ids))
			for i := range res {
				res[i] = ads.Ad{
					ID:    ids[i],
					Title: fmt.Sprintf("test title %d", i),
					Text:  fmt.Sprintf("test text %d", i),
				}
			}
			return res, nil
		})
	return a
}

func notFoundByAliasRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByAlias", mock.Anything, mock.AnythingOfType("string")).
		Return(links.Link{}, ErrNotFound)
	return r
}

func getByAliasRepo(t *testing.T, adIDs []int64) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByAlias", mock.Anything, mock.AnythingOfType("string")).
		Return(links.Link{Ads: adIDs}, nil)
	return r
}

func storeGetByAliasRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("Store", mock.Anything, mock.AnythingOfType("links.Link")).
		Return(func(_ context.Context, link links.Link) (int64, error) {
			if link.URL == "https://already-exists.com" {
				return 0, ErrAlreadyExists
			}
			return 0, nil
		})

	r.
		On("GetByAlias", mock.Anything, mock.AnythingOfType("string")).
		Return(links.Link{}, ErrNotFound)
	return r
}

func getByIDUpdateAliasGetByAliasRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(links.Link{URL: "https://github.com", Alias: "github"}, nil)

	r.
		On("UpdateAlias", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).
		Return(func(_ context.Context, id int64, alias string) error {
			if alias == "already-exists" {
				return ErrAlreadyExists
			}
			return nil
		})

	r.
		On("GetByAlias", mock.Anything, mock.AnythingOfType("string")).
		Return(links.Link{}, ErrNotFound)
	return r
}

func getByIDUpdateAliasRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(links.Link{URL: "https://github.com", Alias: "github"}, nil)

	r.
		On("UpdateAlias", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).
		Return(func(_ context.Context, id int64, alias string) error {
			if alias == "already-exists" {
				return ErrAlreadyExists
			}
			return nil
		})

	return r
}

func getByIDRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(func(_ context.Context, link int64) (links.Link, error) {
			if link == -1 {
				return links.Link{}, ErrNotFound
			}
			return links.Link{URL: "https://github.com", Alias: "github", Ads: []int64{1, 2, 3}}, nil
		})
	return r
}

func generator(t *testing.T) Generator {
	g := mocks.NewGenerator(t)
	g.
		On("Generate", mock.Anything).
		Return("test", nil)
	return g
}

func getByIDAddAdRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(links.Link{Ads: []int64{1, 2, 3}, URL: "https://github.com", Alias: "github"}, nil)
	r.
		On("AddAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(func(_ context.Context, link int64, ad int64) error {
			if ad == -1 {
				return ErrNotFound
			}
			return nil
		})
	return r
}

func getByIDDeleteAdRepo(t *testing.T) Repository {
	r := mocks.NewRepository(t)
	r.
		On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
		Return(links.Link{Ads: []int64{1, 2, 3}, URL: "https://github.com", Alias: "github"}, nil)
	r.
		On("DeleteAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(func(_ context.Context, link int64, ad int64) error {
			if ad == -1 || link == -1 {
				return ErrNotFound
			}
			return nil
		})
	return r
}

func TestApp_AddAd(t *testing.T) {
	type fields struct {
		repo Repository
		gen  Generator
	}
	type args struct {
		ctx      context.Context
		linkID   int64
		adID     int64
		authorID int64
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    links.Link
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct adding",
			fields: fields{
				repo: getByIDAddAdRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   0,
				adID:     0,
				authorID: 0,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3, 0},
			},
		},
		{
			name: "permission denied",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   0,
				adID:     0,
				authorID: 1,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrPermissionDenied, i)
			},
		},
		{
			name: "already been added",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   0,
				adID:     1,
				authorID: 0,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAdAlreadyAdded, i)
			},
		},
		{
			name: "not found",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   -1,
				adID:     1,
				authorID: 0,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo: tt.fields.repo,
				Gen:  tt.fields.gen,
			}
			got, err := a.AddAd(tt.args.ctx, tt.args.linkID, tt.args.adID, tt.args.authorID)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("AddAd(%v, %v, %v, %v)", tt.args.ctx, tt.args.linkID, tt.args.adID, tt.args.authorID)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "AddAd(%v, %v, %v, %v)", tt.args.ctx, tt.args.linkID, tt.args.adID, tt.args.authorID)
		})
	}
}

func TestApp_Create(t *testing.T) {
	type fields struct {
		repo Repository
		gen  Generator
	}
	type args struct {
		ctx      context.Context
		url      string
		alias    string
		authorID int64
		ads      []int64
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    links.Link
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct creating with generating",
			fields: fields{
				repo: storeGetByAliasRepo(t),
				gen:  generator(t),
			},
			args: args{
				ctx: context.Background(),
				url: "https://github.com",
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "test",
				AuthorID: 0,
				Ads:      nil,
			},
		},
		{
			name: "correct creating without generating",
			fields: fields{
				repo: storeRepo(t),
			},
			args: args{
				ctx:   context.Background(),
				url:   "https://github.com",
				alias: "my-alias",
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "my-alias",
				AuthorID: 0,
				Ads:      nil,
			},
		},
		{
			name: "already exists",
			fields: fields{
				repo: storeRepo(t),
			},
			args: args{
				ctx:   context.Background(),
				url:   "https://already-exists.com",
				alias: "test",
			},
			want: links.Link{
				ID:       0,
				URL:      "https://already-exists.com",
				Alias:    "test",
				AuthorID: 0,
				Ads:      nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAlreadyExists, i)
			},
		},
		{
			name: "invalid url",
			fields: fields{
				repo: notFoundByAliasRepo(t),
				gen:  generator(t),
			},
			args: args{
				ctx: context.Background(),
				url: "",
			},
			want: links.Link{
				ID:       0,
				URL:      "",
				Alias:    "test",
				AuthorID: 0,
				Ads:      nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrInvalidContent, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo: tt.fields.repo,
				Gen:  tt.fields.gen,
			}
			got, err := a.Create(tt.args.ctx, tt.args.url, tt.args.alias, tt.args.authorID, tt.args.ads)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("Create(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.url, tt.args.alias, tt.args.authorID, tt.args.ads)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "Create(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.url, tt.args.alias, tt.args.authorID, tt.args.ads)
		})
	}
}

func TestApp_DeleteAd(t *testing.T) {
	type fields struct {
		repo Repository
		gen  Generator
	}
	type args struct {
		ctx      context.Context
		linkID   int64
		adID     int64
		authorID int64
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    links.Link
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct deleting",
			fields: fields{
				repo: getByIDDeleteAdRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   0,
				adID:     2,
				authorID: 0,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 3},
			},
		},
		{
			name: "permission denied",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   0,
				adID:     0,
				authorID: 1,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrPermissionDenied, i)
			},
		},
		{
			name: "ad not found",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   0,
				adID:     0,
				authorID: 0,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
		{
			name: "link not found",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   -1,
				adID:     1,
				authorID: 0,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
		{
			name: "error from delete (prevent deleting from list)",
			fields: fields{
				repo: getByIDDeleteAdRepo(t),
			},
			args: args{
				ctx:      context.Background(),
				linkID:   -1,
				adID:     1,
				authorID: 0,
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo: tt.fields.repo,
				Gen:  tt.fields.gen,
			}
			got, err := a.DeleteAd(tt.args.ctx, tt.args.linkID, tt.args.adID, tt.args.authorID)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("DeleteAd(%v, %v, %v, %v)", tt.args.ctx, tt.args.linkID, tt.args.adID, tt.args.authorID)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "DeleteAd(%v, %v, %v, %v)", tt.args.ctx, tt.args.linkID, tt.args.adID, tt.args.authorID)
		})
	}
}

func TestApp_UpdateAlias(t *testing.T) {
	type fields struct {
		repo Repository
		gen  Generator
	}
	type args struct {
		ctx      context.Context
		id       int64
		authorID int64
		alias    string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    links.Link
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct updating with generating",
			fields: fields{
				repo: getByIDUpdateAliasGetByAliasRepo(t),
				gen:  generator(t),
			},
			args: args{
				ctx: context.Background(),
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "test",
				AuthorID: 0,
				Ads:      nil,
			},
		},
		{
			name: "permission denied",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				authorID: 1,
				ctx:      context.Background(),
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      []int64{1, 2, 3},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrPermissionDenied, i)
			},
		},
		{
			name: "not found",
			fields: fields{
				repo: getByIDRepo(t),
			},
			args: args{
				id:  -1,
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound, i)
			},
		},
		{
			name: "correct creating without generating",
			fields: fields{
				repo: getByIDUpdateAliasRepo(t),
			},
			args: args{
				ctx:   context.Background(),
				alias: "my-alias",
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "my-alias",
				AuthorID: 0,
				Ads:      nil,
			},
		},
		{
			name: "already exists",
			fields: fields{
				repo: getByIDUpdateAliasRepo(t),
			},
			args: args{
				ctx:   context.Background(),
				alias: "already-exists",
			},
			want: links.Link{
				ID:       0,
				URL:      "https://github.com",
				Alias:    "github",
				AuthorID: 0,
				Ads:      nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAlreadyExists, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo: tt.fields.repo,
				Gen:  tt.fields.gen,
			}
			got, err := a.UpdateAlias(tt.args.ctx, tt.args.id, tt.args.authorID, tt.args.alias)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("UpdateAlias(%v, %v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.authorID, tt.args.alias)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "UpdateAlias(%v, %v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.authorID, tt.args.alias)
		})
	}
}

func TestApp_GetRedirect(t *testing.T) {
	type fields struct {
		Repo Repository
		Gen  Generator
		Ads  AdsService
	}
	type args struct {
		ctx   context.Context
		alias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    func(assert.TestingT, redirects.Redirect)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct redirecting",
			fields: fields{
				Repo: getByAliasRepo(t, []int64{1, 2, 3}),
				Ads:  adsService(t),
			},
			args: args{
				ctx:   context.Background(),
				alias: "",
			},
			want: func(t assert.TestingT, redirect redirects.Redirect) {
				if redirect.Ad.ID != 1 && redirect.Ad.ID != 3 {
					t.Errorf("incorrect ad ID. Expected 1 or 3, got: %d", redirect.Ad.ID)
				}
			},
			wantErr: nil,
		},
		{
			name: "ad not found",
			fields: fields{
				Repo: getByAliasRepo(t, []int64{}),
				Ads:  adsService(t),
			},
			args: args{
				ctx:   context.Background(),
				alias: "",
			},
			want: func(t assert.TestingT, redirect redirects.Redirect) {
				assert.Empty(t, redirect.Ad)
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				Repo: tt.fields.Repo,
				Gen:  tt.fields.Gen,
				Ads:  tt.fields.Ads,
			}
			got, err := a.GetRedirect(tt.args.ctx, tt.args.alias)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("GetRedirect(%v, %v)", tt.args.ctx, tt.args.alias)) {
				return
			} else if tt.wantErr == nil {
				assert.NoError(t, err)
			}
			tt.want(t, got)
		})
	}
}

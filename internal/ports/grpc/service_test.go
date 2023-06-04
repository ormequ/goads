package grpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goads/internal/app"
	"goads/internal/app/mocks"
	"goads/internal/entities/ads"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestService_ChangeAdStatus(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *ChangeAdStatusRequest
	}

	amGet := mocks.NewAds(t)
	amGet.
		On("ChangeStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	amGet.
		On("GetByID", mock.Anything, mock.Anything).
		Return(ads.Ad{}, nil)

	amChange := mocks.NewAds(t)
	amChange.
		On("ChangeStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(app.ErrNotFound)

	tests := [...]struct {
		name string
		args args
		ads  *mocks.Ads
		err  error
	}{
		{
			name: "Valid request",
			args: args{request: &ChangeAdStatusRequest{}},
			ads:  amGet,
		},
		{
			name: "Invalid request",
			args: args{request: &ChangeAdStatusRequest{}},
			err:  status.Error(codes.NotFound, app.ErrNotFound.Error()),
			ads:  amChange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Service{
				Ads: tt.ads,
			}
			_, err := a.ChangeAdStatus(tt.args.ctx, tt.args.request)
			fmt.Println(err)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
func TestService_UpdateAd(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *UpdateAdRequest
	}

	amGet := mocks.NewAds(t)
	amGet.
		On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	amGet.
		On("GetByID", mock.Anything, mock.Anything).
		Return(ads.Ad{}, nil)

	amChange := mocks.NewAds(t)
	amChange.
		On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(app.ErrNotFound)

	tests := [...]struct {
		name string
		args args
		ads  *mocks.Ads
		err  error
	}{
		{
			name: "Valid request",
			args: args{request: &UpdateAdRequest{}},
			ads:  amGet,
		},
		{
			name: "Invalid request",
			args: args{request: &UpdateAdRequest{}},
			err:  status.Error(codes.NotFound, app.ErrNotFound.Error()),
			ads:  amChange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Service{
				Ads: tt.ads,
			}
			_, err := a.UpdateAd(tt.args.ctx, tt.args.request)
			fmt.Println(err)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

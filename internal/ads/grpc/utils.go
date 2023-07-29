package grpc

import (
	"errors"
	"goads/internal/ads/ads"
	"goads/internal/ads/app"
	"goads/internal/ads/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorStatus(err error) error {
	if err == nil {
		return nil
	}
	code := codes.Internal
	if errors.Is(err, app.ErrNotFound) {
		code = codes.NotFound
	}
	if errors.Is(err, app.ErrPermissionDenied) {
		code = codes.PermissionDenied
	}
	if errors.Is(err, app.ErrInvalidContent) || errors.Is(err, app.ErrInvalidFilter) {
		code = codes.InvalidArgument
	}
	return status.Error(code, err.Error())
}

func adToResponse(ad ads.Ad) *proto.AdResponse {
	return &proto.AdResponse{
		Id:         ad.ID,
		Title:      ad.Title,
		Published:  ad.Published,
		Text:       ad.Text,
		AuthorId:   ad.AuthorID,
		CreateDate: ad.CreateDate.UnixMilli(),
		UpdateDate: ad.UpdateDate.UnixMilli(),
	}
}

func adsListToResponse(list []ads.Ad) *proto.ListAdResponse {
	res := proto.ListAdResponse{
		List: make([]*proto.AdResponse, len(list)),
	}
	for i := range list {
		res.List[i] = adToResponse(list[i])
	}
	return &res
}

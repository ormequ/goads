package grpc

import (
	"errors"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/links"
	"goads/internal/urlshortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorStatus(err error) error {
	if err == nil {
		return nil
	}
	code := codes.Internal
	if errors.Is(err, app.ErrNoAds) {
		code = codes.OK
	}
	if errors.Is(err, app.ErrAlreadyExists) {
		code = codes.AlreadyExists
	}
	if errors.Is(err, app.ErrNotFound) || errors.Is(err, app.ErrAdNotExists) {
		code = codes.NotFound
	}
	return status.Error(code, err.Error())
}

func linkToResponse(link links.Link) *proto.LinkResponse {
	return &proto.LinkResponse{
		Id:       link.ID,
		Url:      link.URL,
		Alias:    link.Alias,
		AuthorID: link.AuthorID,
		Ads:      link.Ads,
	}
}

func listLinkToResponse(list []links.Link) *proto.ListLinkResponse {
	res := proto.ListLinkResponse{
		List: make([]*proto.LinkResponse, len(list)),
	}
	for i := range list {
		res.List[i] = linkToResponse(list[i])
	}
	return &res
}

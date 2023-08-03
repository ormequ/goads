package grpc

import (
	"errors"
	"goads/internal/pkg/errwrap"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/entities/ads"
	"goads/internal/urlshortener/entities/links"
	"goads/internal/urlshortener/entities/redirects"
	"goads/internal/urlshortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorStatus(err error) error {
	if err == nil {
		return nil
	}
	code := codes.Internal
	var wrap *errwrap.Error
	if errors.As(err, &wrap) {
		err = wrap.Unwrap() // hiding error information
	}
	if errors.Is(err, app.ErrNoAds) {
		code = codes.OK
	}
	if errors.Is(err, app.ErrInvalidContent) {
		code = codes.InvalidArgument
	}
	if errors.Is(err, app.ErrAlreadyExists) || errors.Is(err, app.ErrAdAlreadyAdded) {
		code = codes.AlreadyExists
	}
	if errors.Is(err, app.ErrNotFound) || errors.Is(err, app.ErrAdNotExists) {
		code = codes.NotFound
	}
	if errors.Is(err, app.ErrPermissionDenied) {
		code = codes.PermissionDenied
	}
	if code == codes.Internal {
		err = errors.New("internal error")
	}
	return status.Error(code, err.Error())
}

func linkToResponse(link links.Link) *proto.LinkResponse {
	return &proto.LinkResponse{
		Id:       link.ID,
		Url:      link.URL,
		Alias:    link.Alias,
		AuthorId: link.AuthorID,
		Ads:      link.Ads,
	}
}

func adToResponse(ad ads.Ad) *proto.AdResponse {
	return &proto.AdResponse{
		Id:    ad.ID,
		Title: ad.Title,
		Text:  ad.Text,
	}
}

func redirectToResponse(r redirects.Redirect) *proto.RedirectResponse {
	return &proto.RedirectResponse{
		Link: linkToResponse(r.Link),
		Ad:   adToResponse(r.Ad),
	}
}

func listLinkToResponse(list []links.Link) *proto.LinksListResponse {
	res := proto.LinksListResponse{
		List: make([]*proto.LinkResponse, len(list)),
	}
	for i := range list {
		res.List[i] = linkToResponse(list[i])
	}
	return &res
}

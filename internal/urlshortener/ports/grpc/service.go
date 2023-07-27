package grpc

import (
	"context"
	"goads/internal/urlshortener/links"
	"goads/internal/urlshortener/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type App interface {
	Create(ctx context.Context, url string, alias string, authorID int64, ads []int64) (links.Link, error)
	GetByID(ctx context.Context, id int64) (links.Link, error)
	GetByAuthor(ctx context.Context, author int64) ([]links.Link, error)
	GetByAlias(ctx context.Context, alias string) (links.Link, error)
	UpdateAlias(ctx context.Context, id int64, alias string) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	app App
}

func (s Service) Create(ctx context.Context, request *proto.CreateRequest) (*proto.LinkResponse, error) {
	link, err := s.app.Create(ctx, request.Url, request.Alias, request.AuthorID, request.Ads)
	return linkToResponse(link), getErrorStatus(err)
}

func (s Service) GetByID(ctx context.Context, request *proto.GetByIDRequest) (*proto.LinkResponse, error) {
	link, err := s.app.GetByID(ctx, request.Id)
	return linkToResponse(link), getErrorStatus(err)
}

func (s Service) GetByAuthor(ctx context.Context, request *proto.GetByAuthorRequest) (*proto.ListLinkResponse, error) {
	list, err := s.app.GetByAuthor(ctx, request.AuthorID)
	return listLinkToResponse(list), getErrorStatus(err)
}

func (s Service) GetByAlias(ctx context.Context, request *proto.GetByAliasRequest) (*proto.LinkResponse, error) {
	link, err := s.app.GetByAlias(ctx, request.Alias)
	return linkToResponse(link), getErrorStatus(err)
}

func (s Service) UpdateAlias(ctx context.Context, request *proto.UpdateAliasRequest) (*proto.LinkResponse, error) {
	err := s.app.UpdateAlias(ctx, request.Id, request.Alias)
	if err != nil {
		return nil, getErrorStatus(err)
	}
	link, err := s.app.GetByID(ctx, request.Id)
	return linkToResponse(link), getErrorStatus(err)
}

func (s Service) Delete(ctx context.Context, request *proto.UpdateAliasRequest) (*emptypb.Empty, error) {
	err := s.app.Delete(ctx, request.Id)
	return new(emptypb.Empty), getErrorStatus(err)
}

func NewService(a App) Service {
	return Service{a}
}

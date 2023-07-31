package ads

import (
	"context"
	"fmt"
	"goads/internal/ads/proto"
	"goads/internal/pkg/errwrap"
	"goads/internal/urlshortener/app"
	"goads/internal/urlshortener/entities/ads"
)

type Client struct {
	Svc proto.AdServiceClient
}

func (c Client) GetOnlyPublished(ctx context.Context, ids []int64) ([]ads.Ad, error) {
	const op = "ads.GetOnlyPublished"
	adsList, err := c.Svc.GetOnlyPublished(ctx, &proto.AdIDsRequest{Id: ids})
	if err != nil {
		return nil, errwrap.New(err, app.ServiceName, op).WithDetails(fmt.Sprintf("ad ids: %v", ids))
	}
	res := make([]ads.Ad, len(adsList.List))
	for i := range res {
		res[i] = ads.New(adsList.List[i].Title, adsList.List[i].Text)
	}
	return res, nil
}

func New(svc proto.AdServiceClient) Client {
	return Client{svc}
}

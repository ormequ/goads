package responses

import (
	"github.com/gin-gonic/gin"
	adProto "goads/internal/ads/proto"
	"goads/internal/api/ads/responses"
	shProto "goads/internal/urlshortener/proto"
)

type Redirect struct {
	URL string              `json:"url"`
	Ad  *responses.PublicAd `json:"ad"`
}

type Link struct {
	URL      string  `json:"url"`
	Alias    string  `json:"alias"`
	AuthorID int64   `json:"author_id"`
	Ads      []int64 `json:"ads"`
}

func RedirectToResponse(l *shProto.LinkResponse, ad *adProto.AdResponse) Redirect {
	if l == nil {
		return Redirect{}
	}
	if ad == nil {
		return Redirect{URL: l.Url}
	}
	resAd := responses.AdToPublicResponse(ad)
	return Redirect{
		URL: l.Url,
		Ad:  &resAd,
	}
}

func LinkToResponse(l *shProto.LinkResponse) Link {
	if l == nil {
		return Link{}
	}
	return Link{
		URL:      l.Url,
		Alias:    l.Alias,
		AuthorID: l.AuthorId,
		Ads:      l.Ads,
	}
}

func LinksListToResponse(l *shProto.LinksListResponse) []Link {
	if l == nil {
		return nil
	}
	res := make([]Link, len(l.List))
	for i := range res {
		res[i] = LinkToResponse(l.List[i])
	}
	return res
}

func RedirectSuccess(l *shProto.LinkResponse, ad *adProto.AdResponse) gin.H {
	return gin.H{
		"data":  RedirectToResponse(l, ad),
		"error": nil,
	}
}

func LinksListSuccess(l *shProto.LinksListResponse) gin.H {
	return gin.H{
		"data":  LinksListToResponse(l),
		"error": nil,
	}
}

func LinkSuccess(l *shProto.LinkResponse) gin.H {
	return gin.H{
		"data":  LinkToResponse(l),
		"error": nil,
	}
}

func EmptySuccess() gin.H {
	return gin.H{
		"data":  "",
		"error": nil,
	}
}

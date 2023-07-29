package responses

import (
	"github.com/gin-gonic/gin"
	"goads/internal/ads/proto"
	"time"
)

type PublicAd struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Ad struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"author_id"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
	Published  bool      `json:"published"`
	PublicAd   `json:"public"`
}

func AdToPublicResponse(a *proto.AdResponse) PublicAd {
	if a == nil {
		return PublicAd{}
	}
	return PublicAd{
		Title: a.Title,
		Text:  a.Text,
	}
}

func AdToResponse(a *proto.AdResponse) Ad {
	if a == nil {
		return Ad{}
	}
	return Ad{
		ID:         a.Id,
		AuthorID:   a.AuthorId,
		CreateDate: time.UnixMilli(a.CreateDate).UTC(),
		UpdateDate: time.UnixMilli(a.UpdateDate).UTC(),
		Published:  a.Published,
		PublicAd:   AdToPublicResponse(a),
	}
}

func AdSuccess(a *proto.AdResponse) gin.H {
	return gin.H{
		"data":  AdToResponse(a),
		"error": nil,
	}
}

func AdsListSuccess(l *proto.ListAdResponse) gin.H {
	responses := make([]Ad, len(l.List))
	for i, ad := range l.List {
		responses[i] = AdToResponse(ad)
	}
	return gin.H{
		"data":  responses,
		"error": nil,
	}
}

func EmptySuccess() gin.H {
	return gin.H{
		"data":  "",
		"error": nil,
	}
}

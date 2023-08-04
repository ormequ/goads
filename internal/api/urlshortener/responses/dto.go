package responses

import (
	"github.com/gin-gonic/gin"
	"goads/internal/urlshortener/proto"
)

type Ad struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Redirect struct {
	URL string `json:"url"`
	Ad  *Ad    `json:"ad"`
}

type Link struct {
	URL      string  `json:"url"`
	Alias    string  `json:"alias"`
	AuthorID int64   `json:"author_id"`
	Ads      []int64 `json:"ads"`
}

func RedirectToResponse(r *proto.RedirectResponse) Redirect {
	if r == nil {
		return Redirect{}
	}
	if r.Ad == nil || r.Ad.Title == "" && r.Ad.Text == "" {
		return Redirect{URL: r.Link.Url}
	}
	resAd := Ad{
		Title: r.Ad.Title,
		Text:  r.Ad.Text,
	}
	return Redirect{
		URL: r.Link.Url,
		Ad:  &resAd,
	}
}

func LinkToResponse(l *proto.LinkResponse) Link {
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

func LinksToResponse(l *proto.LinksResponse) []Link {
	if l == nil {
		return nil
	}
	res := make([]Link, len(l.List))
	for i := range res {
		res[i] = LinkToResponse(l.List[i])
	}
	return res
}

func RedirectSuccess(r *proto.RedirectResponse) gin.H {
	return gin.H{
		"data":  RedirectToResponse(r),
		"error": nil,
	}
}

func LinksSuccess(l *proto.LinksResponse) gin.H {
	return gin.H{
		"data":  LinksToResponse(l),
		"error": nil,
	}
}

func LinkSuccess(l *proto.LinkResponse) gin.H {
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

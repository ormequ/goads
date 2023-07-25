package responses

import (
	"github.com/gin-gonic/gin"
	"goads/internal/ads/ads"
	"time"
)

type ad struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"author_id"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
	Published  bool      `json:"published"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
}

func adToResponse(a ads.Ad) ad {
	return ad{
		ID:         a.ID,
		AuthorID:   a.AuthorID,
		CreateDate: a.CreateDate,
		UpdateDate: a.UpdateDate,
		Published:  a.Published,
		Title:      a.Title,
		Text:       a.Text,
	}
}

func AdSuccess(ad ads.Ad) gin.H {
	return gin.H{
		"data":  adToResponse(ad),
		"error": nil,
	}
}

func AdsListSuccess(adsList []ads.Ad) gin.H {
	responses := make([]ad, len(adsList))
	for i, ad := range adsList {
		responses[i] = adToResponse(ad)
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

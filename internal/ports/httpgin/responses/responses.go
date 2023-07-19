package responses

import (
	"github.com/gin-gonic/gin"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
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

type user struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

func userToResponse(u users.User) user {
	return user{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
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

func UserSuccess(user users.User) gin.H {
	return gin.H{
		"data":  userToResponse(user),
		"error": nil,
	}
}

func EmptySuccess() gin.H {
	return gin.H{
		"data":  "",
		"error": nil,
	}
}

func Error(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

package httpgin

import (
	"github.com/gin-gonic/gin"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title" binding:"required"`
	Text   string `json:"text" binding:"required"`
	UserID int64  `json:"user_id"`
}

type createUserRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

type adResponse struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"author_id"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
	Published  bool      `json:"published"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type deleteAdRequest struct {
	UserID int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title" binding:"required"`
	Text   string `json:"text" binding:"required"`
	UserID int64  `json:"user_id"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type changeUserNameRequest struct {
	Name string `json:"name" binding:"required"`
}
type changeUserEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func adToResponse(ad ads.Ad) adResponse {
	return adResponse{
		ID:         ad.ID,
		AuthorID:   ad.AuthorID,
		CreateDate: ad.CreateDate,
		UpdateDate: ad.UpdateDate,
		Published:  ad.Published,
		Title:      ad.Title,
		Text:       ad.Text,
	}
}

func userToResponse(user users.User) userResponse {
	return userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func AdSuccessResponse(ad ads.Ad) gin.H {
	return gin.H{
		"data":  adToResponse(ad),
		"error": nil,
	}
}

func AdsListSuccessResponse(adsList []ads.Ad) gin.H {
	responses := make([]adResponse, len(adsList))
	for i, ad := range adsList {
		responses[i] = adToResponse(ad)
	}
	return gin.H{
		"data":  responses,
		"error": nil,
	}
}

func UserSuccessResponse(user users.User) gin.H {
	return gin.H{
		"data":  userToResponse(user),
		"error": nil,
	}
}

func EmptySuccessResponse() gin.H {
	return gin.H{
		"data":  "",
		"error": nil,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

package responses

import (
	"github.com/gin-gonic/gin"
	"goads/internal/auth/proto"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UserToResponse(u *proto.UserInfoResponse) User {
	if u == nil {
		return User{}
	}
	return User{
		ID:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}

func TokenToResponse(t *proto.TokenResponse) string {
	if t == nil {
		return ""
	}
	return t.Token
}

func UserSuccess(u *proto.UserInfoResponse) gin.H {
	return gin.H{
		"data":  UserToResponse(u),
		"error": nil,
	}
}

func TokenSuccess(t *proto.TokenResponse) gin.H {
	return gin.H{
		"data":  TokenToResponse(t),
		"error": nil,
	}
}

func UserWithTokenSuccess(reg *proto.RegisterResponse) gin.H {
	if reg == nil {
		return gin.H{
			"data":  nil,
			"error": nil,
		}
	}
	return gin.H{
		"data": gin.H{
			"User":  UserToResponse(reg.User),
			"token": TokenToResponse(reg.Token),
		},
		"error": nil,
	}
}

func EmptySuccess() gin.H {
	return gin.H{
		"data":  "",
		"error": nil,
	}
}

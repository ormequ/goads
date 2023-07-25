package responses

import (
	"github.com/gin-gonic/gin"
	"goads/internal/auth/users"
)

type user struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func userToResponse(u users.User) user {
	return user{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func UserSuccess(user users.User) gin.H {
	return gin.H{
		"data":  userToResponse(user),
		"error": nil,
	}
}

func TokenSuccess(token string) gin.H {
	return gin.H{
		"data":  token,
		"error": nil,
	}
}

func UserWithTokenSuccess(user users.User, token string) gin.H {
	return gin.H{
		"data": gin.H{
			"user":  userToResponse(user),
			"token": token,
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

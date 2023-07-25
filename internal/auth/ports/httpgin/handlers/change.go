package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/auth/ports/httpgin/responses"
	"goads/internal/auth/ports/httpgin/utils"
	"net/http"
)

type emailChanger interface {
	getterByID
	validator
	ChangeEmail(ctx context.Context, id int64, email string) error
}

type changeUserEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func ChangeEmail(app emailChanger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserEmailRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := app.Validate(c, utils.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.Error(err))
			return
		}
		err = app.ChangeEmail(c, user.ID, req.Email)
		user.Email = req.Email
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

type nameChanger interface {
	getterByID
	validator
	ChangeName(ctx context.Context, id int64, name string) error
}

type changeUserNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func ChangeName(app nameChanger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserNameRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := app.Validate(c, utils.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.Error(err))
			return
		}
		err = app.ChangeName(c, user.ID, req.Name)
		user.Name = req.Name
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

type passwordChanger interface {
	getterByID
	validator
	ChangePassword(ctx context.Context, id int64, password string) error
}

type changeUserPasswordRequest struct {
	Password string `json:"name" binding:"required"`
}

func ChangePassword(app passwordChanger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserPasswordRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		user, err := app.Validate(c, utils.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.Error(err))
			return
		}
		err = app.ChangePassword(c, user.ID, req.Password)
		if err != nil {
			displayHiddenErr(c, err)
			return
		}
		c.JSON(http.StatusOK, responses.UserSuccess(user))
	}
}

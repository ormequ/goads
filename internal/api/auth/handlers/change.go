package handlers

import (
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/responses"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/auth/proto"
	"net/http"
)

type changeEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func ChangeEmail(client proto.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeEmailRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		id, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		user, err := client.ChangeEmail(c, &proto.ChangeUserEmailRequest{
			Id:    id,
			Email: req.Email,
		})
		errors.ProceedResult(c, responses.UserSuccess(user), err)
	}
}

type changeNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func ChangeName(app proto.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeNameRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		id, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		user, err := app.ChangeName(c, &proto.ChangeUserNameRequest{
			Id:   id,
			Name: req.Name,
		})
		errors.ProceedResult(c, responses.UserSuccess(user), err)
	}
}

type changePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func ChangePassword(app proto.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changePasswordRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		id, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		user, err := app.ChangePassword(c, &proto.ChangeUserPasswordRequest{
			Id:       id,
			Password: req.Password,
		})
		errors.ProceedResult(c, responses.UserSuccess(user), err)
	}
}

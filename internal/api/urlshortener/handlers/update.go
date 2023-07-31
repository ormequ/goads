package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/api/urlshortener/responses"
	"goads/internal/urlshortener/proto"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type updateAdRequest struct {
	Ad int64 `json:"ad" binding:"required"`
}

func UpdateAdData(method func(context.Context, *proto.LinkAdRequest, ...grpc.CallOption) (*proto.LinkResponse, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("link_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
		}
		var req updateAdRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		userID, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		link, err := method(c, &proto.LinkAdRequest{
			LinkId:   int64(id),
			AdId:     req.Ad,
			AuthorId: userID,
		})
		errors.ProceedResult(c, responses.LinkSuccess(link), err)
	}
}

type updateAliasRequest struct {
	Alias string `json:"alias" binding:"required"`
}

func UpdateAlias(shortener proto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("link_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
		}
		var req updateAliasRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
			return
		}
		userID, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		link, err := shortener.UpdateAlias(c, &proto.UpdateAliasRequest{
			Id:       int64(id),
			AuthorId: userID,
			Alias:    req.Alias,
		})
		errors.ProceedResult(c, responses.LinkSuccess(link), err)
	}
}

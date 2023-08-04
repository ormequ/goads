package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/api/urlshortener/responses"
	"goads/internal/urlshortener/proto"
	"net/http"
	"strconv"
)

func GetRedirect(shortener proto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("alias")
		redirect, err := shortener.GetRedirect(c, &proto.GetByAliasRequest{Alias: alias})
		errors.ProceedResult(c, responses.RedirectSuccess(redirect), err)
	}
}

func GetByID(shortener proto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("link_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.Response(err))
		}
		userID, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		link, err := shortener.GetByID(c, &proto.GetByIDRequest{Id: int64(id)})
		if link.AuthorId != userID {
			c.JSON(http.StatusForbidden, errors.Response(fmt.Errorf("access denied")))
		}
		errors.ProceedResult(c, responses.LinkSuccess(link), err)
	}
}

func GetByAuthor(shortener proto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		links, err := shortener.GetByAuthor(c, &proto.GetByAuthorRequest{AuthorId: id})
		errors.ProceedResult(c, responses.LinksSuccess(links), err)
	}
}

package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	adProto "goads/internal/ads/proto"
	"goads/internal/api/auth/utils"
	"goads/internal/api/errors"
	"goads/internal/api/urlshortener/responses"
	shProto "goads/internal/urlshortener/proto"
	"math/rand"
	"net/http"
	"strconv"
)

func GetByAlias(shortener shProto.ShortenerServiceClient, ads adProto.AdServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("alias")
		link, err := shortener.GetByAlias(c, &shProto.GetByAliasRequest{Alias: alias})
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		adsList, err := ads.GetOnlyPublished(c, &adProto.AdIDsRequest{Id: link.Ads})
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		var ad *adProto.AdResponse
		if len(adsList.List) > 0 {
			ad = adsList.List[rand.Intn(len(adsList.List))]
		} else {
			ad = nil
		}
		c.JSON(http.StatusOK, responses.RedirectSuccess(link, ad))
	}
}

func GetByID(shortener shProto.ShortenerServiceClient) gin.HandlerFunc {
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
		link, err := shortener.GetByID(c, &shProto.GetByIDRequest{Id: int64(id)})
		if link.AuthorId != userID {
			c.JSON(http.StatusForbidden, errors.Response(fmt.Errorf("access denied")))
		}
		errors.ProceedResult(c, responses.LinkSuccess(link), err)
	}
}

func GetByAuthor(shortener shProto.ShortenerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.GetUserID(c)
		if err != nil {
			c.JSON(errors.GetHTTPStatus(err), errors.HiddenResponse(err))
			return
		}
		links, err := shortener.GetByAuthor(c, &shProto.GetByAuthorRequest{AuthorId: id})
		errors.ProceedResult(c, responses.LinksListSuccess(links), err)
	}
}

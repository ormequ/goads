package httpgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goads/internal/app"
	"goads/internal/entities/ads"
	"goads/internal/entities/users"
	"goads/internal/filters"
	"net/http"
	"strconv"
	"time"
)

func createAd(a app.Ads) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createAdRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		ad, err := a.Create(c, req.Title, req.Text, req.UserID)
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		fmt.Println(ad)
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func changeAdStatus(a app.Ads) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeAdStatusRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		err = a.ChangeStatus(c, int64(id), req.UserID, req.Published)
		var ad ads.Ad
		if err == nil {
			ad, err = a.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func updateAd(a app.Ads) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateAdRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		err = a.Update(c, int64(id), req.UserID, req.Title, req.Text)
		var ad ads.Ad
		if err == nil {
			ad, err = a.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAdsFiltered(a app.Ads) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := filters.AdsOptions{
			AuthorID: -1,
		}
		for _, key := range [3]string{"all", "date", "author"} {
			val, ok := c.GetQuery(key)
			if !ok {
				continue
			}
			switch key {
			case "all":
				filter.All = true
			case "author":
				id, err := strconv.Atoi(val)
				if err != nil {
					c.JSON(http.StatusBadRequest, ErrorResponse(err))
					return
				}
				filter.AuthorID = int64(id)
			case "date":
				date, err := time.Parse(time.RFC3339Nano, val)
				if err != nil {
					c.JSON(http.StatusBadRequest, ErrorResponse(err))
					return
				}
				filter.Date = date
			}
		}
		adsList, err := a.GetFiltered(c, filter)
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdsListSuccessResponse(adsList))
	}
}

func searchAds(a app.Ads) gin.HandlerFunc {
	return func(c *gin.Context) {
		adsList, err := a.Search(c, c.Query("q"))
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdsListSuccessResponse(adsList))
	}
}

func deleteAd(a app.Ads) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req deleteAdRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		err = a.Delete(c, int64(id), req.UserID)
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, EmptySuccessResponse())
	}
}

func createUser(a app.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		user, err := a.Create(c, req.Email, req.Name)
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func getUser(a app.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		user, err := a.GetByID(c, int64(id))
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func changeUserName(a app.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserNameRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		err = a.ChangeName(c, int64(id), req.Name)
		var user users.User
		if err == nil {
			user, err = a.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func changeUserEmail(a app.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changeUserEmailRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		err = a.ChangeEmail(c, int64(id), req.Email)
		var user users.User
		if err == nil {
			user, err = a.GetByID(c, int64(id))
		}
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func deleteUser(a app.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		err = a.Delete(c, int64(id))
		if err != nil {
			c.JSON(getErrorHTTPStatus(err), ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, EmptySuccessResponse())
	}
}

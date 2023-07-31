package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (int64, error) {
	s, get := c.Get("userID")
	id, conv := s.(int64)
	if !get {
		return id, errors.New("unable to read user ID from context")
	}
	if !conv {
		return id, fmt.Errorf("unable to convert ID %s to int64", s)
	}
	return id, nil
}

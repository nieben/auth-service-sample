package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/api"
	"github.com/nieben/auth-service-sample/model"
	"net/http"
	"time"
)

var (
	TokenRequiredErr = errors.New("token required")
	TokenInvalidErr  = errors.New("invalid token")
	TokenExpiredErr  = errors.New("token expired")
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.NewFailResponse(TokenRequiredErr, nil))
			return
		}
		if len(token) != model.TokenLength {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.NewFailResponse(TokenInvalidErr, nil))
			return
		}

		t := model.GetToken(token)
		if t == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.NewFailResponse(TokenInvalidErr, nil))
			return
		}
		if t.ExpireAt < time.Now().Unix() {
			t.Remove()
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.NewFailResponse(TokenExpiredErr, nil))
			return
		}
		if model.GetUser(t.User.Username) == nil { // user deleted
			t.Remove()
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.NewFailResponse(model.UserNotExistErr, nil))
			return
		}

		c.Set("user", t.User)
		c.Set("token", t)
		c.Next()
	}
}

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/api"
	"github.com/nieben/auth-service-sample/model"
	"net/http"
)

type AuthController struct {
}

// @Summary token
// @Tags auth
// @Accept json
// @Produce json
// @Param data body api.Token true "请求参数"
// @Success 200 {object} api.Response{}
// @Header  200 {string} Token ""
// @Router /auth/token [post]
func (a *AuthController) Token(c *gin.Context) {
	var in api.Token
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := in.Check(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	user := model.GetUser(in.Username)
	if user == nil {
		c.JSON(http.StatusOK, api.NewFailResponse(model.UserCheckErr, nil))
		return
	}
	if user.Status == 1 { // deleted
		c.JSON(http.StatusOK, api.NewFailResponse(model.UserNotExistErr, nil))
		return
	}
	if !user.CheckPwd(in.Password) {
		c.JSON(http.StatusOK, api.NewFailResponse(model.UserCheckErr, nil))
		return
	}

	t := model.GenerateToken(user)
	c.Header("token", t.Token)
	c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
}

// @Summary logout
// @Tags auth
// @Accept json
// @Produce json
// @Param token header string true "请求参数"
// @Success 200 {object} api.Response{}
// @Router /auth/logout [post]
func (a *AuthController) Logout(c *gin.Context) {
	t, _ := c.Get("token")
	if err := t.(*model.Token).Remove(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
	} else {
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

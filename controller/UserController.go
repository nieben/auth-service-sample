package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/api"
	"github.com/nieben/auth-service-sample/model"
	"net/http"
)

type UserController struct {
}

// @Summary create user
// @Tags user
// @Accept json
// @Produce json
// @Param data body api.CreateUser true "请求参数"
// @Success 200 {object} api.Response{}
// @Router /user/create [post]
func (u *UserController) Create(c *gin.Context) {
	var in api.CreateUser
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := in.Check(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := model.CreateUser(in.Username, in.Password); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	} else {
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

// @Summary delete user
// @Tags user
// @Accept json
// @Produce json
// @Param data body api.DeleteUser true "请求参数"
// @Success 200 {object} api.Response{}
// @Router /user/delete [post]
func (u *UserController) Delete(c *gin.Context) {
	var in api.DeleteUser
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := in.Check(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := model.DeleteUser(in.Username); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	} else {
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

// @Summary add role to user
// @Tags user
// @Accept json
// @Produce json
// @Param data body api.AddUserRole true "请求参数"
// @Success 200 {object} api.Response{}
// @Router /user/addRole [post]
func (u *UserController) AddRole(c *gin.Context) {
	var in api.AddUserRole
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
		c.JSON(http.StatusOK, api.NewFailResponse(model.UserNotExistErr, nil))
		return
	}
	role := model.GetRole(in.Role)
	if role == nil {
		c.JSON(http.StatusOK, api.NewFailResponse(model.RoleNotExistErr, nil))
		return
	}

	if err := user.AddRole(role); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	} else {
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

// @Summary check role
// @Tags user
// @Accept json
// @Produce json
// @Param token header string true "请求参数"
// @Param data body api.CheckRole true "请求参数"
// @Success 200 {object} api.Response{data=bool}
// @Router /user/checkRole [post]
func (u *UserController) CheckRole(c *gin.Context) {
	var in api.CheckRole
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := in.Check(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	role := model.GetRole(in.Role)
	if role == nil {
		c.JSON(http.StatusOK, api.NewFailResponse(model.RoleNotExistErr, nil))
		return
	}

	user, _ := c.Get("user")
	c.JSON(http.StatusOK, api.NewSuccessResponse(user.(*model.User).CheckRole(in.Role)))
}

// @Summary roles
// @Tags user
// @Accept json
// @Produce json
// @Param token header string true "请求参数"
// @Success 200 {object} api.Response{data=[]string}
// @Router /user/roles [post]
func (u *UserController) Roles(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, api.NewSuccessResponse(user.(*model.User).Roles()))
}

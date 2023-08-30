package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/api"
	"github.com/nieben/auth-service-sample/model"
	"net/http"
)

type RoleController struct {
}

// @Summary create role
// @Tags role
// @Accept json
// @Produce json
// @Param data body api.CreateRole true "请求参数"
// @Success 200 {object} api.Response{}
// @Router /role/create [post]
func (r *RoleController) Create(c *gin.Context) {
	var in api.CreateRole
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := in.Check(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := model.CreateRole(in.Role); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
	} else {
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

// @Summary delete role
// @Tags role
// @Accept json
// @Produce json
// @Param data body api.DeleteRole true "请求参数"
// @Success 200 {object} api.Response{}
// @Router /role/delete [post]
func (r *RoleController) Delete(c *gin.Context) {
	var in api.DeleteRole
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := in.Check(); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
		return
	}

	if err := model.DeleteRole(in.Role); err != nil {
		c.JSON(http.StatusOK, api.NewFailResponse(err, nil))
	} else {
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

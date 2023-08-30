package route

import (
	"github.com/gin-gonic/gin"
	"github.com/nieben/auth-service-sample/controller"
	"github.com/nieben/auth-service-sample/route/middleware"
)

var (
	userController = &controller.UserController{}
	roleController = &controller.RoleController{}
	authController = &controller.AuthController{}
)

func Init() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.PrintInfo())

	user := router.Group("/user")
	{
		user.POST("/create", userController.Create)
		user.POST("/delete", userController.Delete)
		user.POST("/addRole", userController.AddRole)
		user.POST("/checkRole", middleware.TokenAuth(), userController.CheckRole)
		user.POST("/roles", middleware.TokenAuth(), userController.Roles)
	}

	role := router.Group("/role")
	{
		role.POST("/create", roleController.Create)
		role.POST("/delete", roleController.Delete)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/token", authController.Token)
		auth.POST("/logout", middleware.TokenAuth(), authController.Logout)
	}

	return router
}

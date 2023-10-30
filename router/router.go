package router

import (
	"rakamin-final/controllers"
	"rakamin-final/middlewares"

	"github.com/gin-gonic/gin"
)

func ConfigureRouter() *gin.Engine {
	router := gin.Default()

	publicRoutes := router.Group("/api")
	{
		publicRoutes.POST("/users/login", controllers.HandleLogin)
		publicRoutes.POST("/users/register", controllers.HandleUserRegistration)
	}

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middlewares.Authenticate())
	{
		protectedRoutes.GET("/users/:id", controllers.HandleUserByID)
		protectedRoutes.PUT("/users/:id", controllers.HandleUpdateUser)
		protectedRoutes.DELETE("/users/:id", controllers.HandleDeleteUser)
		protectedRoutes.GET("/photos", controllers.HandleAllPhotos)
		protectedRoutes.GET("/photos/:id", controllers.HandlePhotoByID)
		protectedRoutes.POST("/photos", controllers.HandleNewPhoto)
	}

	return router
}

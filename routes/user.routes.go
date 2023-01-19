package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/controllers"
	"github.com/wpcodevo/google-github-oath2-golang/middleware"
)

func UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser())
	router.GET("/me", controllers.GetMe)
}

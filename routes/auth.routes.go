package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/controllers"
	"github.com/wpcodevo/google-github-oath2-golang/middleware"
)

func AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", controllers.SignUpUser)
	router.POST("/login", controllers.SignInUser)
	router.GET("/logout", middleware.DeserializeUser(), controllers.LogoutUser)
}

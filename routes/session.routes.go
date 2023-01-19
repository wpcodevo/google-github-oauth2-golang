package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/controllers"
)

func SessionRoute(rg *gin.RouterGroup) {
	router := rg.Group("/sessions/oauth")

	router.GET("/google", controllers.GoogleOAuth)
	router.GET("/github", controllers.GitHubOAuth)
}

package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
	"github.com/wpcodevo/google-github-oath2-golang/models"
	"github.com/wpcodevo/google-github-oath2-golang/utils"
)

func SignUpUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  payload.Password,
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(&newUser)}})
}

func SignInUser(ctx *gin.Context) {
	var payload *models.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if user.Provider == "Google" || user.Provider == "GitHub" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": fmt.Sprintf("Use %v OAuth instead", user.Provider)})
		return
	}

	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := utils.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	google_user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	email := strings.ToLower(google_user.Email)

	user_data := models.User{
		Name:      google_user.Name,
		Email:     email,
		Password:  "",
		Photo:     google_user.Picture,
		Provider:  "Google",
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if initializers.DB.Model(&user_data).Where("email = ?", email).Updates(&user_data).RowsAffected == 0 {
		initializers.DB.Create(&user_data)
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", email)

	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID.String(), config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathUrl))
}

package controllers

import (
	"ApiGin/dtos"
	"ApiGin/initializers"
	"ApiGin/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(cgin *gin.Context) {

	var body struct {
		Email    string
		Password string
		Name     string
		IsAdmin  bool
	}

	if err := cgin.ShouldBindJSON(&body); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 15)
	if err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create the User

	user := models.User{Email: body.Email, Password: string(hash), Name: body.Name, IsAdmin: body.IsAdmin}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		//Failed to create the user
		cgin.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	cgin.JSON(http.StatusCreated, gin.H{"result": "User Created"})

}

func Login(cgin *gin.Context) {
	// var minutes,_ = strconv.ParseInt(os.Getenv("EXPIRATION_TIME"), 2, 64)
	// min, _ := strconv.ParseInt(minutes, 2, 64)

	var expiration time.Duration = 240
	var body struct {
		Email    string
		Password string
	}

	if err := cgin.ShouldBindJSON(&body); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.Id == 0 {
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Invalid Email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		// Invalid password
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Invalid Email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Id,
		"exp":   time.Now().Add(time.Minute * expiration).Unix(),
		"email": user.Email,
		"name":  user.Name,
		"admin": user.IsAdmin,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		// Failed to generate the token
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Failed to generate the token"})
		return
	}

	var response dtos.UserLogin
	response.Token = tokenString
	response.Expiration = time.Now().Add(time.Minute * expiration)
	response.Email = user.Email
	response.Admin = user.IsAdmin
	cgin.JSON(http.StatusOK, response)
}

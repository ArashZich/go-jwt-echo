package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arashzich/echo-go/initializers"
	"github.com/arashzich/echo-go/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {

	// Get data off req body
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid Request")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Create a user
	user := models.User{
		Username: body.Username,
		Password: string(hashedPassword),
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// return it
	return c.JSON(http.StatusOK, user)
}

func Login(c echo.Context) error {

	// Get the user_name and password
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get the user
	var user models.User
	result := initializers.DB.Where("username = ?", body.Username).First(&user)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found or password is wrong"})
	}

	// Compare the password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong Password"})
	}

	// Create a token and refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 20).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// return it
	return c.JSON(http.StatusOK, map[string]string{
		"accessToken":  tokenString,
		"refreshToken": refreshTokenString,
	})
}

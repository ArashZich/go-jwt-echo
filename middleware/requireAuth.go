// package middleware

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/arashzich/echo-go/initializers"
// 	"github.com/arashzich/echo-go/models"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/labstack/echo/v4"
// )

// func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		// Get the token from headers
// 		tokenString := c.Request().Header.Get("Authorization")

// 		// Check if token is empty
// 		if tokenString == "" {
// 			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token is required"})
// 		}

// 		// Parse the token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Don't forget to validate the alg is what you expect:
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

// 			}

// 			// hmacSampleSecret is a []byte containing your secret, e.g., []byte("my_secret_key")
// 			return []byte(os.Getenv("SECRET")), nil
// 		})

// 		if err != nil {
// 			log.Println("Error parsing token:", err)
// 			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			// Check the expiration time
// 			if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token has expired"})
// 			}

// 			// Find the user with token sub
// 			var user models.User
// 			if err := initializers.DB.First(&user, claims["sub"]).Error; err != nil {
// 				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
// 			}

// 			// Attach user to the context
// 			c.Set("user", user)

// 			// Continue with the next middleware or the actual handler
// 			return next(c)
// 		}

// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
// 	}
// }

package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/arashzich/echo-go/initializers"
	"github.com/arashzich/echo-go/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token is required"})
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
		}

		// Extract the token without the Bearer prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Secret key for HMAC validation
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		// Check token validity and claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the expiration time
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token has expired"})
			}

			// Find the user with token sub
			var user models.User
			if err := initializers.DB.First(&user, claims["sub"]).Error; err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
			}

			// Attach user to the context
			c.Set("user", user)

			// Continue with the next middleware or the actual handler
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}
}

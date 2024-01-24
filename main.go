package main

import (
	"github.com/labstack/echo/v4"

	"github.com/arashzich/echo-go/controllers"
	"github.com/arashzich/echo-go/initializers"
	"github.com/arashzich/echo-go/middleware" // Import the RequireAuth middleware
)

func init() {
	initializers.LoadEbvVariables()
	initializers.ConnectToDb()

}

func main() {

	app := echo.New()

	// Unprotected routes
	app.POST("/signup", controllers.Signup)
	app.POST("/login", controllers.Login)

	// Group the /posts routes under RequireAuth middleware
	postsGroup := app.Group("/posts")
	postsGroup.Use(middleware.RequireAuth) // Apply RequireAuth middleware to /posts routes

	// Protected routes under /posts
	postsGroup.POST("", controllers.PostCreate)
	postsGroup.GET("", controllers.PostIndex)
	postsGroup.GET("/:id", controllers.PostShow)
	postsGroup.PUT("/:id", controllers.PostUpdate)
	postsGroup.DELETE("/:id", controllers.PostDelete)

	app.Logger.Fatal(app.Start(":3000"))

}

package main

import (
	"github.com/arashzich/echo-go/initializers"
	"github.com/arashzich/echo-go/models"
)

func init() {
	initializers.LoadEbvVariables()
	initializers.ConnectToDb()

}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Post{})

}

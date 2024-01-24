package controllers

import (
	"net/http"

	"github.com/arashzich/echo-go/initializers"
	"github.com/arashzich/echo-go/models"
	"github.com/labstack/echo/v4"
)

func PostCreate(c echo.Context) error {
	// Get data off req body
	var body struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	c.Bind(&body)

	// Create a post
	post := models.Post{
		Title: body.Title,
		Body:  body.Body,
	}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// return it

	return c.JSON(http.StatusOK, post)
}

func PostIndex(c echo.Context) error {

	// Get The post
	var posts []models.Post
	result := initializers.DB.Find(&posts)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// return it
	return c.JSON(http.StatusOK, posts)
}

func PostShow(c echo.Context) error {
	// Get the id
	id := c.Param("id")

	// Get the post
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// return it
	return c.JSON(http.StatusOK, post)
}

func PostUpdate(c echo.Context) error {
	// Get the id
	id := c.Param("id")

	// Get the post
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// Get data off req body
	var body struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	c.Bind(&body)

	// Update the post
	post.Title = body.Title
	post.Body = body.Body

	result = initializers.DB.Save(&post)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// return it
	return c.JSON(http.StatusOK, post)
}

func PostDelete(c echo.Context) error {

	// Get the id
	id := c.Param("id")

	// Get the post
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// Delete the post
	result = initializers.DB.Delete(&post)

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, result.Error)
	}

	// return it
	return c.JSON(http.StatusOK, "Post Deleted")

}

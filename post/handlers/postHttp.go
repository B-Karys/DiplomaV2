package handlers

import (
	"DiplomaV2/post/models"
	"DiplomaV2/post/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type postHttpHandler struct {
	postUsecase usecase.PostUseCase
}

// HTTP handler
func (p postHttpHandler) CreatePost(c echo.Context) error {
	//Extract user ID from JWT token or any other authentication mechanism
	userID := c.Get("userID").(int64)

	// Parse request body
	var input struct { //future: grpc message
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Convert the string to int64

	//Create a new Post object
	post := &models.Post{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		AuthorID:    userID,
	}

	// Call the use case to create a new post
	err := p.postUsecase.CreatePost(post)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Successfully created post"})
}

func (p postHttpHandler) UpdatePost(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (p postHttpHandler) DeletePost(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (p postHttpHandler) GetPost(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewPostHttpHandler(postUsecase usecase.PostUseCase) PostHandler {
	return &postHttpHandler{
		postUsecase: postUsecase,
	}
}

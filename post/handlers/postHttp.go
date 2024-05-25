package handlers

import (
	"DiplomaV2/post/models"
	"DiplomaV2/post/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type postHttpHandler struct {
	postUsecase usecase.PostUseCase
}

func (p *postHttpHandler) CreatePost(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var input struct { //future: grpc message
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Skills      []string `json:"skills"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	post := &models.Post{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Skills:      input.Skills,
		AuthorID:    userID,
	}

	err := p.postUsecase.CreatePost(post)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Successfully created post"})
}

func (p *postHttpHandler) UpdatePost(c echo.Context) error {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	postIDs := c.Param("id")
	postID, err := strconv.ParseInt(postIDs, 10, 64)

	authorID := c.Get("userID").(int64)

	err = p.postUsecase.UpdatePost(postID, authorID, input.Name, input.Description, input.Type)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post updated successfully"})
}

func (p *postHttpHandler) DeletePost(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (p *postHttpHandler) GetPost(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewPostHttpHandler(postUsecase usecase.PostUseCase) PostHandler {
	return &postHttpHandler{
		postUsecase: postUsecase,
	}
}

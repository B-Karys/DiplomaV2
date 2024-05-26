package handlers

import (
	"DiplomaV2/internal/helpers"
	"DiplomaV2/internal/validator"
	"DiplomaV2/post/models"
	"DiplomaV2/post/usecase"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
)

type postHttpHandler struct {
	postUseCase usecase.PostUseCase
}

func (p *postHttpHandler) GetMyPosts(c echo.Context) error {
	userID := c.Get("userID").(int64)
	userIDstring := strconv.FormatInt(userID, 10)
	var input struct {
		Name        string
		Description string
		PostType    string
		Skills      []string
		models.Filters
	}

	v := validator.New()

	qs := c.Request().URL.Query()

	//Search by words
	input.Name = helpers.ReadString(qs, "name", "")
	input.Description = helpers.ReadString(qs, "description", "")

	//Show only one type
	input.PostType = helpers.ReadString(qs, "type", "")

	//Show the posts of certain author: will be needed for certain author profiles
	//Show only by chosen skills
	input.Skills = helpers.ReadCSV(qs, "skills", []string{})

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "pageSize", 20, v)
	input.Filters.Sort = helpers.ReadString(qs, "sort", "created_at")
	input.Filters.SortSafeList = []string{"name", "created_at", "-name", "-created_at"}

	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	posts, err := p.postUseCase.GetFilteredPosts(input.Name, input.Description, userIDstring, input.PostType, input.Skills, input.Filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

func (p *postHttpHandler) GetFilteredPosts(c echo.Context) error {
	var input struct {
		Name        string
		Description string
		Author      string
		PostType    string
		Skills      []string
		models.Filters
	}

	v := validator.New()

	qs := c.Request().URL.Query()

	//Search by words
	input.Name = helpers.ReadString(qs, "name", "")
	input.Description = helpers.ReadString(qs, "description", "")

	//Show only one type
	input.PostType = helpers.ReadString(qs, "type", "")

	//Show the posts of certain author: will be needed for certain author profiles
	input.Author = helpers.ReadString(qs, "author", "")

	//Show only by chosen skills
	input.Skills = helpers.ReadCSV(qs, "skills", []string{})

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "pageSize", 20, v)
	input.Filters.Sort = helpers.ReadString(qs, "sort", "created_at")
	input.Filters.SortSafeList = []string{"name", "created_at", "-name", "-created_at"}

	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	posts, err := p.postUseCase.GetFilteredPosts(input.Name, input.Description, input.Author, input.PostType, input.Skills, input.Filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

func (p *postHttpHandler) DeletePost(c echo.Context) error {
	userID := c.Get("userID")
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post id"})
	}

	post, err := p.postUseCase.GetPostById(postID)
	if post == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Post not found"})
	}

	if post.AuthorID != userID || err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Post doesn't belong to you"})
	}
	err = p.postUseCase.DeletePost(post.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (p *postHttpHandler) GetPostById(c echo.Context) error {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post id"})
	}
	post, err := p.postUseCase.GetPostById(postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, post)
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
		Type:        strings.ToLower(input.Type),
		Skills:      input.Skills,
		AuthorID:    userID,
	}

	err := p.postUseCase.CreatePost(post)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Successfully created post"})
}

func (p *postHttpHandler) UpdatePost(c echo.Context) error {
	var input struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Type        string         `json:"type"`
		Skills      pq.StringArray `json:"skills"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	postIDs := c.Param("id")
	postID, err := strconv.ParseInt(postIDs, 10, 64)

	authorID := c.Get("userID").(int64)

	err = p.postUseCase.UpdatePost(postID, authorID, input.Name, input.Description, input.Skills, strings.ToLower(input.Type))
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post updated successfully"})
}

func NewPostHttpHandler(postUseCase usecase.PostUseCase) PostHandler {
	return &postHttpHandler{
		postUseCase: postUseCase,
	}
}

package handlers

import (
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/internal/helpers"
	"DiplomaV2/backend/internal/validator"
	postsFilter "DiplomaV2/backend/post"
	"DiplomaV2/backend/post/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type postHttpHandler struct {
	postUseCase usecase.PostUseCase
}

func (p *postHttpHandler) GetMyPosts(c echo.Context) error {
	userID := c.Get("userID").(int64)
	userIDString := strconv.FormatInt(userID, 10)
	var input struct {
		Name        string
		Description string
		Author      string
		PostType    string
		Skills      []string
		postsFilter.Filters
	}

	v := validator.New()

	qs := c.Request().URL.Query()

	input.Name = helpers.ReadString(qs, "name", "")
	input.Description = helpers.ReadString(qs, "description", "")
	input.PostType = helpers.ReadString(qs, "type", "")
	input.Author = helpers.ReadString(qs, "author", userIDString)
	input.Skills = helpers.ReadCSV(qs, "skills", []string{})

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "pageSize", 10, v)
	input.Filters.Sort = helpers.ReadString(qs, "sort", "created_at")
	input.Filters.SortSafeList = []string{"name", "created_at", "-name", "-created_at"}

	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	if postsFilter.ValidateFilters(v, input.Filters); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	post := entity.Post{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.PostType,
		AuthorID:    userID,
		Skills:      input.Skills,
	}

	posts, metadata, err := p.postUseCase.GetFilteredPosts(&post, input.Filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	type Response struct {
		Posts    []*entity.Post       `json:"posts"`
		Metadata postsFilter.Metadata `json:"metadata"`
	}

	response := Response{
		Posts:    posts,
		Metadata: metadata,
	}

	return c.JSON(http.StatusOK, response)
}

func (p *postHttpHandler) GetFilteredPosts(c echo.Context) error {
	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		AuthorID    int64    `json:"author_id"`
		PostType    string   `json:"type"`
		Skills      []string `json:"skills"`
		postsFilter.Filters
	}

	v := validator.New()

	qs := c.Request().URL.Query()

	input.Name = helpers.ReadString(qs, "name", "")
	input.Description = helpers.ReadString(qs, "description", "")
	input.PostType = helpers.ReadString(qs, "type", "")
	input.AuthorID = int64(helpers.ReadInt(qs, "author", 0, v))
	input.Skills = helpers.ReadCSV(qs, "skills", []string{})

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "pageSize", 10, v)
	input.Filters.Sort = helpers.ReadString(qs, "sort", "created_at")
	input.Filters.SortSafeList = []string{"name", "created_at", "-name", "-created_at"}

	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	if postsFilter.ValidateFilters(v, input.Filters); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	post := entity.Post{
		Name:        input.Name,
		Description: input.Description,
		AuthorID:    input.AuthorID,
		Type:        input.PostType,
		Skills:      input.Skills,
	}

	posts, metadata, err := p.postUseCase.GetFilteredPosts(&post, input.Filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	type Response struct {
		Posts    []*entity.Post       `json:"posts"`
		Metadata postsFilter.Metadata `json:"metadata"`
	}

	response := Response{
		Posts:    posts,
		Metadata: metadata,
	}

	return c.JSON(http.StatusOK, response)
}

func (p *postHttpHandler) DeletePost(c echo.Context) error {
	userID := c.Get("userID").(int64)
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

	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		PostType    string   `json:"type"`
		Skills      []string `json:"skills"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	post := entity.Post{
		Name:        input.Name,
		Description: input.Description,
		Type:        strings.ToLower(input.PostType),
		Skills:      input.Skills,
		AuthorID:    userID,
	}

	err := p.postUseCase.CreatePost(&post)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Successfully created post"})
}

func (p *postHttpHandler) UpdatePost(c echo.Context) error {
	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		PostType    string   `json:"type"`
		Skills      []string `json:"skills"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	postIDs := c.Param("id")
	postID, err := strconv.ParseInt(postIDs, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post id"})
	}

	authorID := c.Get("userID").(int64)

	post := entity.Post{
		Name:        input.Name,
		Description: input.Description,
		Type:        strings.ToLower(input.PostType),
		Skills:      input.Skills,
	}

	err = p.postUseCase.UpdatePost(postID, authorID, &post)
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

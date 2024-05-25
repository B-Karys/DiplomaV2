package handlers

import (
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

func (p *postHttpHandler) GetFilteredPosts(c echo.Context) error {
	// Parse query parameters to extract filters
	authorID := c.QueryParam("author_id")
	postType := c.QueryParam("type")

	// Convert postType to lowercase for searching
	postType = strings.ToLower(postType)

	// Call the appropriate use case method based on the presence of filters
	if authorID != "" || postType != "" {
		// If filters are provided, call the GetFilteredPosts method
		posts, err := p.postUseCase.GetFilteredPosts(authorID, postType)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, posts)
	} else {
		// If no filters are provided, call the GetAllPosts method
		posts, err := p.postUseCase.GetAllPosts()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, posts)
	}
}

func (p *postHttpHandler) GetMyPosts(c echo.Context) error {
	userID := c.Get("userID").(int64)
	// Parse query parameters to extract filters
	authorID := strconv.FormatInt(userID, 10)
	println("My User ID is: ", authorID)

	postType := c.QueryParam("type")

	// Call the appropriate use case method based on the presence of filters
	if authorID != "" || postType != "" {
		// If filters are provided, call the GetFilteredPosts method
		posts, err := p.postUseCase.GetFilteredPosts(authorID, postType)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, posts)
	} else {
		// If no filters are provided, call the GetAllPosts method
		posts, err := p.postUseCase.GetAllPosts()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, posts)
	}
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

/*
func (p *postHttpHandler) GetPostByAuthor(c echo.Context) error {
	var input struct {
		AuthorID int64 `json:"author_id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	posts, err := p.postUseCase.GetPostByAuthor(input.AuthorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, posts)
}
*/

//func (p *postHttpHandler) GetAllPosts(c echo.Context) error {
//	// Extract filter criteria from query parameters
//	authorID := c.QueryParam("author_id")
//	postType := c.QueryParam("type")
//	// Assuming more filters like skills, etc.
//
//	// Call the use case method passing the filter criteria
//	posts, err := p.postUseCase.GetFilteredPosts(authorID, postType /*, additional filter criteria*/)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
//	}
//	return c.JSON(http.StatusOK, posts)
//}

/*
func (p *postHttpHandler) GetPostByType(c echo.Context) error {

	posts, err := p.postUseCase.GetPostByType(c.QueryParam("type"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, posts)
}
*/

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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post updated successfully"})
}

func NewPostHttpHandler(postUseCase usecase.PostUseCase) PostHandler {
	return &postHttpHandler{
		postUseCase: postUseCase,
	}
}

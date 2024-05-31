package handlers

import "github.com/labstack/echo/v4"

type PostHandler interface {
	CreatePost(c echo.Context) error
	GetPostById(c echo.Context) error
	GetFilteredPosts(c echo.Context) error // Changed method name
	GetMyPosts(c echo.Context) error
	UpdatePost(c echo.Context) error
	DeletePost(c echo.Context) error
}

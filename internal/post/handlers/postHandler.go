package handlers

import "github.com/labstack/echo/v4"

type PostHandler interface {
	CreatePost(c echo.Context) error
	updatePost(c echo.Context) error
	deletePost(c echo.Context) error
	getPost(c echo.Context) error
}

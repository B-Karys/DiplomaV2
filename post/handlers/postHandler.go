package handlers

import "github.com/labstack/echo/v4"

type PostHandler interface {
	CreatePost(c echo.Context) error
	UpdatePost(c echo.Context) error
	DeletePost(c echo.Context) error
	GetPost(c echo.Context) error
}

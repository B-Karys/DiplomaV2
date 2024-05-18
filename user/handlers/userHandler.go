package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	Registration(c echo.Context) error
	Activation(c echo.Context) error
	GetUserInfoByEmail(c echo.Context) error
	GetUserInfoById(c echo.Context) error
	UpdateUserInfo(c echo.Context) error
	DeleteUser(c echo.Context) error
	ResetPassword(c echo.Context) error
}

package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	Registration(c echo.Context) error
	Activation(c echo.Context) error
	Authentication(c echo.Context) error
	CheckAuth(c echo.Context) error
	GetUserInfoByEmail(c echo.Context) error
	GetUserInfoById(c echo.Context) error
	UpdateUserInfo(c echo.Context) error
	ChangePassword(c echo.Context) error
	ResetPassword(c echo.Context) error
	ForgotPassword(c echo.Context) error
	Logout(c echo.Context) error
	DeleteUser(c echo.Context) error
}

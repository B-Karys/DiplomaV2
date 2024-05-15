package handlers

import (
	"github.com/labstack/echo/v4"
)

type TokenHandler interface {
	CreateActivationToken(c echo.Context) error
	CreateAuthenticationToken(c echo.Context) error
	CreatePasswordResetToken(c echo.Context) error
}

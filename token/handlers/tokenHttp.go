package handlers

import (
	"DiplomaV2/internal/validator"
	"DiplomaV2/token/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrFailedValidation = errors.New("Failed Validation")
)

type tokenHttpHandler struct {
	tokenUseCase usecase.TokenUseCase
}

func (t tokenHttpHandler) CreateActivationToken(c echo.Context) error {

	panic("implement me")
}

func (t tokenHttpHandler) CreateAuthenticationToken(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	v := validator.New()
	validator.ValidateEmail(v, input.Email)
	validator.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": ErrFailedValidation.Error()})
	}

	token, err := t.tokenUseCase.CreateAuthenticationToken(input.Email, input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"authentication_token": token})
}

func (t tokenHttpHandler) CreatePasswordResetToken(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewTokenHttpHandler(tokenUseCase usecase.TokenUseCase) TokenHandler {
	return &tokenHttpHandler{
		tokenUseCase: tokenUseCase,
	}
}

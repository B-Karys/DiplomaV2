package handlers

import (
	"DiplomaV2/internal/validator"
	"DiplomaV2/user/models"
	"DiplomaV2/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrFailedValidation = errors.New("Failed Validation")
)

type userHttpHandler struct {
	userUsecase usecase.UserUseCase
}

func (u userHttpHandler) Activation(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHttpHandler) ResetPassword(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHttpHandler) Registration(c echo.Context) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := &models.User{
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		Activated: true,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		println("Error setting password")
		return err
	}

	v := validator.New()
	if validator.ValidateUser(v, user); !v.Valid() {
		return ErrFailedValidation
	}

	err = u.userUsecase.Registration(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = c.JSON(http.StatusCreated, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return err
}

func (u userHttpHandler) GetUserInfo(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHttpHandler) UpdateUserInfo(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHttpHandler) DeleteUser(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewUserHttpHandler(userUsecase usecase.UserUseCase) UserHandler {
	return &userHttpHandler{userUsecase}
}

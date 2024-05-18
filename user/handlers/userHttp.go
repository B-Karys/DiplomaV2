package handlers

import (
	"DiplomaV2/internal/mailer"
	"DiplomaV2/internal/validator"
	"DiplomaV2/user/models"
	"DiplomaV2/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrFailedValidation = errors.New("Failed Validation")
	ErrWrongCredentials = errors.New("Wrong Credentials")
	ErrWrongPassword    = errors.New("Wrong Password")
)

type userHttpHandler struct {
	userUsecase usecase.UserUseCase
	mailer      mailer.Mailer
}

func (u userHttpHandler) Authentication(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	}

	v := validator.New()
	validator.ValidateEmail(v, input.Email)
	validator.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		return c.JSON(http.StatusBadRequest, ErrWrongCredentials.Error())
	}

	user, err := u.userUsecase.GetUserByEmail(input.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrWrongCredentials.Error())
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrWrongPassword.Error())
	}

	if !match {
		return c.JSON(http.StatusBadRequest, ErrWrongPassword.Error())
	}

	token, err := u.userUsecase.Authentication(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrWrongCredentials.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"authentication_token": token})
}

func (u userHttpHandler) Activation(c echo.Context) error {
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	}

	v := validator.New()
	if validator.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	}

	err := u.userUsecase.Activation(input.TokenPlaintext)
	if err != nil {
		println("activation error: " + err.Error())
		return c.JSON(http.StatusInternalServerError, ErrFailedValidation.Error())
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{"message": "Activation successful"})
}

func (u userHttpHandler) ResetPassword(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHttpHandler) ChangePassword(c echo.Context) error {
	/*
		var input struct {
			currentPassword string `json:"currentPassword"`
			newPassword     string `json:"newPassword"`
			repeatNewPass   string `json:"repeatNewPass"`
		}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, ErrFailedValidation)
		}
		v := validator.New()
		if validator.ValidateTokenPlaintext(v, input.currentPassword); !v.Valid() {
			return c.JSON(http.StatusBadRequest, ErrFailedValidation)
		}
		//TODO implement me with middleware
	*/
	return nil
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
		Activated: false,
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

	token, err := u.userUsecase.Registration(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	u.background(func() error {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		return u.mailer.Send(user.Email, "user_welcome.tmpl", data)
	})

	err = c.JSON(http.StatusCreated, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return err
}

func (u userHttpHandler) background(fn func() error) {
	// Launch a background goroutine.
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		err := fn()
		if err != nil {
			return
		}
	}()
}

func (u userHttpHandler) GetUserInfoByEmail(c echo.Context) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := u.userUsecase.GetUserByEmail(input.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (u userHttpHandler) GetUserInfoById(c echo.Context) error {
	var input struct {
		Id int64 `json:"id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := u.userUsecase.GetUserById(input.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (u userHttpHandler) UpdateUserInfo(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHttpHandler) DeleteUser(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewUserHttpHandler(userUsecase usecase.UserUseCase, theMailer mailer.Mailer) UserHandler {
	return &userHttpHandler{userUsecase,
		theMailer,
	}
}

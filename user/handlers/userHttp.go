package handlers

import (
	"DiplomaV2/internal/mailer"
	"DiplomaV2/internal/validator"
	middleware2 "DiplomaV2/middleware"
	"DiplomaV2/user/models"
	"DiplomaV2/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrFailedValidation = errors.New("validation error")
	ErrWrongCredentials = errors.New("wrong Credentials")
	ErrWrongPassword    = errors.New("wrong Password")
	ErrNotActive        = errors.New("user is not activated")
	ErrNotValidPassword = errors.New("password is not valid")
)

type userHttpHandler struct {
	userUseCase usecase.UserUseCase
	mailer      mailer.Mailer
}

func (u *userHttpHandler) Authentication(c echo.Context) error {
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

	user, err := u.userUseCase.GetUserByEmail(input.Email)
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

	if user.Activated != true {
		return c.JSON(http.StatusForbidden, ErrNotActive.Error())
	}

	token, err := u.userUseCase.Authentication(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrWrongCredentials.Error())
	}

	expiry := time.Now().Add(24 * time.Hour)

	c.SetCookie(&http.Cookie{
		Expires:  expiry,
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Secure:   true, // Рекомендуется использовать только с HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	// Set the user in the context
	c.Set("user", user)

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User successfully authenticated"})
}

func (u *userHttpHandler) CheckAuth(c echo.Context) error {
	// Use LoginMiddleware to extract and validate the JWT token from the cookie
	err := middleware2.LoginMiddleware(func(c echo.Context) error {
		// Token is valid, user is authenticated
		return c.JSON(http.StatusOK, map[string]string{"authenticated": "true"})
	})(c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"authenticated": "false"})
	}

	return nil
}

func (u *userHttpHandler) Activation(c echo.Context) error {
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

	err := u.userUseCase.Activation(input.TokenPlaintext)
	if err != nil {
		println("activation error: " + err.Error())
		return c.JSON(http.StatusInternalServerError, ErrFailedValidation.Error())
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{"message": "Activation successful"})
}

func (u *userHttpHandler) ResetPassword(c echo.Context) error {
	//var input struct {
	//	TokenPlaintext string `json:"token"`
	//	NewPassword    string `json:"new_password"`
	//}
	//if err := c.Bind(&input); err != nil {
	//	return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	//}
	//v := validator.New()
	//if validator.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
	//	return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	//}
	//
	//if validator.ValidatePasswordPlaintext(v, input.NewPassword); !v.Valid() {
	//	return c.JSON(http.StatusBadRequest, ErrNotValidPassword)
	//}
	//u.userUseCase.ResetPassword()
	//TODO with the second handler
	return nil
}

func (u *userHttpHandler) ChangePassword(c echo.Context) error {
	var input struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
		RepeatNewPass   string `json:"repeatNewPass"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	}

	if input.NewPassword != input.RepeatNewPass {
		return c.JSON(http.StatusBadRequest, "New passwords do not match")
	}

	userId := c.Get("userID").(int64)

	err := u.userUseCase.ChangePassword(userId, input.CurrentPassword, input.NewPassword)
	if err != nil {
		if errors.Is(err, ErrWrongPassword) {
			return c.JSON(http.StatusBadRequest, "Current password is incorrect")
		}
		return c.JSON(http.StatusInternalServerError, "Failed to update password")
	}

	return c.JSON(http.StatusOK, "Password updated successfully")
}

func (u *userHttpHandler) Registration(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	}

	token, err := u.userUseCase.Registration(user)
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

	// If no error, send the response with status accepted and user data
	return c.JSON(http.StatusCreated, map[string]interface{}{"user": user})
}

func (u *userHttpHandler) background(fn func() error) {
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

func (u *userHttpHandler) GetUserInfoByEmail(c echo.Context) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := u.userUseCase.GetUserByEmail(input.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (u *userHttpHandler) GetUserInfoById(c echo.Context) error {
	// Extract the id parameter from the URL
	idParam := c.Param("id")

	// Convert the id parameter from string to int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// Call the use case to get user information by id
	user, err := u.userUseCase.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	data := &models.User{
		Email:    user.Email,
		Name:     user.Name,
		Surname:  user.Surname,
		Username: user.Username,
		Telegram: user.Telegram,
		Discord:  user.Discord,
		Skills:   user.Skills,
	}

	// Return the user information
	return c.JSON(http.StatusOK, data)
}

func (u *userHttpHandler) UpdateUserInfo(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u *userHttpHandler) DeleteUser(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u *userHttpHandler) Logout(c echo.Context) error {
	/* //TODO: once the refresh tokens are done
	token := c.Get("context_token").(string)
	println("context token is:", token)

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
	}

	err := u.userUseCase.DeleteToken(userID, token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	*/

	// Clear the authentication cookie
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Logout successful"})
}

func NewUserHttpHandler(userUsecase usecase.UserUseCase, theMailer mailer.Mailer) UserHandler {
	return &userHttpHandler{userUsecase,
		theMailer,
	}
}

package handlers

import (
	"DiplomaV2/backend/internal/helpers"
	"DiplomaV2/backend/internal/mailer"
	middleware2 "DiplomaV2/backend/internal/middleware"
	"DiplomaV2/backend/internal/validator"
	"DiplomaV2/backend/user/models"
	"DiplomaV2/backend/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrFailedValidation         = errors.New("validation error")
	ErrPasswordFailedValidation = errors.New("Password validation error")
	ErrWrongCredentials         = errors.New("wrong Credentials")
	ErrWrongPassword            = errors.New("wrong Password")
	ErrNotActive                = errors.New("user is not activated")
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

func (u *userHttpHandler) GetMyInfo(c echo.Context) error {
	// Get the user ID from the context
	userID := c.Get("userID").(int64)

	// Retrieve the user's profile data from the database
	user, err := u.userUseCase.GetUserById(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return the user's profile data
	return c.JSON(http.StatusOK, user)
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
	tokenPlaintext := c.Param("token")

	v := validator.New()
	if validator.ValidateTokenPlaintext(v, tokenPlaintext); !v.Valid() {
		return c.JSON(http.StatusBadRequest, ErrFailedValidation.Error())
	}

	err := u.userUseCase.Activation(tokenPlaintext)
	if err != nil {
		println("activation error: " + err.Error())
		return c.JSON(http.StatusInternalServerError, ErrFailedValidation.Error())
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{"message": "Activation successful"})
}

func (u *userHttpHandler) ForgotPassword(c echo.Context) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Call the use case to generate the token
	token, err := u.userUseCase.ForgotPassword(input.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	forgotPasswordLink := fmt.Sprintf("http://localhost:5173/reset-password/%s", token)

	// Send the password reset email in the background
	u.background(func() error {
		data := map[string]any{
			"passwordResetToken": token,
			"forgotPasswordLink": forgotPasswordLink,
		}
		return u.mailer.Send(input.Email, "token_password_reset.tmpl", data)
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset email sent"})
	//Todo with frontend
}

/*
token, err := u.userUseCase.Registration(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	activationLink := fmt.Sprintf("http://localhost:4000/v2/users/activate/%s", token.Plaintext)

	u.background(func() error {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
			"activationLink":  activationLink,
		}
		return u.mailer.Send(user.Email, "user_welcome.tmpl", data)
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{"user": user})
}
*/

func (u *userHttpHandler) ResetPassword(c echo.Context) error {
	var input struct {
		Token           string `json:"token"`
		NewPassword     string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	v := validator.New()
	if validator.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	if input.NewPassword != input.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, ErrPasswordFailedValidation.Error())
	}

	if validator.ValidatePasswordPlaintext(v, input.NewPassword); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	err := u.userUseCase.ResetPassword(input.Token, input.NewPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successful"})
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

	activationLink := fmt.Sprintf("http://localhost:4000/v2/users/activate/%s", token.Plaintext)

	u.background(func() error {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
			"activationLink":  activationLink,
		}
		return u.mailer.Send(user.Email, "user_welcome.tmpl", data)
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{"user": user})
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
	// Return the user information
	return c.JSON(http.StatusOK, user)
}
func (u *userHttpHandler) UpdateUserInfo(c echo.Context) error {
	// Log the incoming request
	fmt.Println("Updating user info...")

	// Parse the multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Error parsing form data:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse form data"})
	}

	// Extract user data
	name := form.Value["name"][0]
	surname := form.Value["surname"][0]
	telegram := form.Value["telegram"][0]
	discord := form.Value["discord"][0]
	skills := form.Value["skills"]
	fmt.Println("Received data:", name, surname, telegram, discord, skills)

	// Get the user ID from the context
	userID := c.Get("userID").(int64)

	// Handle profile image upload
	profileImage := form.File["profileImage"]
	var profileImageURL string
	if len(profileImage) > 0 {
		file := profileImage[0]
		fmt.Println("Received file:", file.Filename)

		src, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
		}
		defer src.Close()

		// Initialize GCS client
		ctx := context.Background()
		client, err := helpers.NewStorageClient(ctx, "C:/Users/krump/Downloads/lucid-volt-424719-f0-5df86076a210.json")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initialize GCS client"})
		}

		// Upload the file to GCS
		objectName := fmt.Sprintf("users/%d/profileImage", userID) // Unique object name based on user ID
		if err := helpers.UploadFileToGCS(ctx, client, "teamfinderimages", objectName, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload file to GCS"})
		}

		// Construct the profile image URL
		profileImageURL = fmt.Sprintf("https://storage.googleapis.com/teamfinderimages/users/%d/profileImage", userID)
	}

	// Call the use case to update the user
	err = u.userUseCase.UpdateUserInfo(userID, name, surname, telegram, discord, skills, profileImageURL)
	if err != nil {
		fmt.Println("Error updating user info:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("User updated successfully")
	return c.JSON(http.StatusAccepted, map[string]string{"message": "User updated successfully"})
}

func (u *userHttpHandler) DeleteUser(c echo.Context) error {

	contextUser := c.Get("userID").(int64)

	// Call the repository method to delete the user
	err := u.userUseCase.DeleteUser(contextUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	return c.NoContent(http.StatusNoContent)
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

func NewUserHttpHandler(userUsecase usecase.UserUseCase, theMailer mailer.Mailer) UserHandler {
	return &userHttpHandler{userUsecase,
		theMailer,
	}
}

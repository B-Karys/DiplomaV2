package handlers

import (
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/internal/helpers"
	"DiplomaV2/backend/internal/mailer"
	middleware2 "DiplomaV2/backend/internal/middleware"
	"DiplomaV2/backend/internal/validator"
	"DiplomaV2/backend/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
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
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.Set("user", user)

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User successfully authenticated"})
}

func (u *userHttpHandler) GetMyInfo(c echo.Context) error {
	userID := c.Get("userID").(int64)

	user, err := u.userUseCase.GetUserById(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (u *userHttpHandler) CheckAuth(c echo.Context) error {
	err := middleware2.LoginMiddleware(func(c echo.Context) error {
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

	token, err := u.userUseCase.ForgotPassword(input.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	forgotPasswordLink := fmt.Sprintf("http://localhost:5173/reset-password/%s", token)

	u.background(func() error {
		data := map[string]any{
			"passwordResetToken": token,
			"forgotPasswordLink": forgotPasswordLink,
		}
		return u.mailer.Send(input.Email, "token_password_reset.tmpl", data)
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset email sent"})
}

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

func (u *userHttpHandler) GetAllUsers(c echo.Context) error {
	type UserInfo struct {
		ID       int64          `json:"id"`
		Name     string         `json:"name"`
		Surname  string         `json:"surname"`
		Username string         `json:"username"`
		Telegram string         `json:"telegram"`
		Discord  string         `json:"discord"`
		Skills   pq.StringArray `json:"skills"`
	}

	users, err := u.userUseCase.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var usersInfo []UserInfo
	for _, user := range users {
		userInfo := UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Surname:  user.Surname,
			Username: user.Username,
			Telegram: user.Telegram,
			Discord:  user.Discord,
			Skills:   user.Skills,
		}
		usersInfo = append(usersInfo, userInfo)
	}

	return c.JSON(http.StatusOK, usersInfo)
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

	user := &entity.User{
		Name:         input.Name,
		Username:     input.Username,
		Email:        input.Email,
		ProfileImage: "https://storage.googleapis.com/teamfinderimages/default_photo.png",
		Activated:    false,
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
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := u.userUseCase.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (u *userHttpHandler) UpdateUserInfo(c echo.Context) error {
	fmt.Println("Updating user info...")

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Error parsing form data:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse form data"})
	}

	name := form.Value["name"][0]
	surname := form.Value["surname"][0]
	username := form.Value["username"][0]
	telegram := form.Value["telegram"][0]
	discord := form.Value["discord"][0]
	skills := form.Value["skills"]
	fmt.Println(form.Value["name"], form.Value["username"])
	fmt.Println("Received data:", name, surname, username, telegram, discord, skills)

	userID := c.Get("userID").(int64)

	profileImage := form.File["profileImage"]
	var profileImageURL string
	if len(profileImage) > 0 {
		file := profileImage[0]
		fmt.Println("Received file:", file.Filename)

		profileImageURL, err = u.userUseCase.UploadProfileImage(userID, file)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	err = u.userUseCase.UpdateUserInfo(userID, name, surname, username, telegram, discord, skills, profileImageURL)
	if err != nil {
		fmt.Println("Error updating user info:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("User updated successfully")
	return c.JSON(http.StatusAccepted, map[string]string{"message": "User updated successfully"})
}

func (u *userHttpHandler) DeleteUser(c echo.Context) error {
	// Get the userID
	userID := c.Get("userID").(int64)

	// Retrieve user information to get the profile image URL
	userInfo, err := u.userUseCase.GetUserById(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user information"})
	}

	// If profile image URL exists, delete the corresponding object from GCS
	if userInfo.ProfileImage != "" {
		// Extract object name from the profile image URL
		objectName := userInfo.ProfileImage[strings.LastIndex(userInfo.ProfileImage, "/")+1:]

		// Initialize GCS client
		ctx := context.Background()
		client, err := helpers.NewStorageClient(ctx, "C:/Users/krump/Downloads/lucid-volt-424719-f0-5df86076a210.json")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initialize GCS client"})
		}

		// Delete the object from GCS
		if err := helpers.DeleteFileFromGCS(ctx, client, "teamfinderimages", objectName); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete profile image from GCS"})
		}
	}

	// Call the repository method to delete the user
	if err := u.userUseCase.DeleteUser(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (u *userHttpHandler) Logout(c echo.Context) error {
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

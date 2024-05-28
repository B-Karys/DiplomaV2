package server

import (
	"DiplomaV2/backend/internal/config"
	"DiplomaV2/backend/internal/database"
	"DiplomaV2/backend/internal/mailer"
	middleware2 "DiplomaV2/backend/internal/middleware"
	postHandlers "DiplomaV2/backend/post/handlers"
	postModels "DiplomaV2/backend/post/models"
	postRepositories "DiplomaV2/backend/post/repository"
	postUseCases "DiplomaV2/backend/post/usecase"
	userHandlers "DiplomaV2/backend/user/handlers"
	userModels "DiplomaV2/backend/user/models"
	userRepositories "DiplomaV2/backend/user/repository"
	tokenRepositories "DiplomaV2/backend/user/tokenRepository"
	userUseCases "DiplomaV2/backend/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"path/filepath"
)

type echoServer struct {
	app    *echo.Echo
	db     database.Database
	conf   *config.Config
	mailer mailer.Mailer
}

func NewEchoServer(conf *config.Config, db database.Database) Server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)
	appMailer := mailer.New("sandbox.smtp.mailtrap.io", 25, "6f71a6ef2443f6", "57f94aefae5b38", "Test <no-reply@test.com>")

	return &echoServer{
		app:    echoApp,
		db:     db,
		conf:   conf,
		mailer: appMailer,
	}
}

func (s *echoServer) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	s.app.Static("/uploads", filepath.Join(os.Getenv("HOME"), "Desktop", "uploads"))
	// CORS middleware with configuration
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"}, // Specify your React frontend domain here
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true, // Allow credentials (cookies)

	}))

	// Handle OPTIONS requests
	s.app.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// Health route
	s.app.GET("v2/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	s.initializeMigrations()

	// Initialize Handlers
	s.initializePostHttpHandler()
	s.initializeUserHttpHandler()

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *echoServer) initializeMigrations() {
	err := s.db.GetDb().AutoMigrate(
		&userModels.User{},
		&postModels.Post{},
		&userModels.Token{},
	)
	if err != nil {
		return
	}
}

func (s *echoServer) initializeUserHttpHandler() {
	userPostgresRepository := userRepositories.NewUserRepository(s.db)
	tokenPostgresRepository := tokenRepositories.NewTokenRepository(s.db)
	userUseCase := userUseCases.NewUserUseCase(userPostgresRepository, tokenPostgresRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(userUseCase, s.mailer)

	userRouters := s.app.Group("/v2/users")
	{
		userRouters.POST("/registration", userHttpHandler.Registration)
		userRouters.GET("/activate/:token", userHttpHandler.Activation)
		userRouters.POST("/login", userHttpHandler.Authentication)
		userRouters.GET("/check-auth", userHttpHandler.CheckAuth)
		userRouters.GET("/:id", userHttpHandler.GetUserInfoById, middleware2.LoginMiddleware)
		userRouters.GET("/my", userHttpHandler.GetMyInfo, middleware2.LoginMiddleware)
		userRouters.PATCH("/update", userHttpHandler.UpdateUserInfo, middleware2.LoginMiddleware)
		userRouters.PATCH("/password", userHttpHandler.ChangePassword, middleware2.LoginMiddleware) //
		userRouters.POST("/logout", userHttpHandler.Logout, middleware2.LoginMiddleware)
		userRouters.DELETE("/:id", userHttpHandler.DeleteUser, middleware2.LoginMiddleware) //
		userRouters.POST("/forgot-password", userHttpHandler.ForgotPassword)                //
		userRouters.POST("/reset-password", userHttpHandler.ResetPassword)                  //
	}
}

func (s *echoServer) initializePostHttpHandler() {
	postPostgresRepository := postRepositories.NewPostRepository(s.db)
	postUseCase := postUseCases.NewPostUseCase(postPostgresRepository)
	postHttpHandler := postHandlers.NewPostHttpHandler(postUseCase)

	postRouters := s.app.Group("/v2/posts")
	postRouters.POST("/", postHttpHandler.CreatePost, middleware2.LoginMiddleware)
	postRouters.GET("/:id", postHttpHandler.GetPostById)
	postRouters.GET("/", postHttpHandler.GetFilteredPosts)
	postRouters.GET("/my", postHttpHandler.GetMyPosts, middleware2.LoginMiddleware)
	postRouters.PATCH("/:id", postHttpHandler.UpdatePost, middleware2.LoginMiddleware)
	postRouters.DELETE("/:id", postHttpHandler.DeletePost, middleware2.LoginMiddleware)

}

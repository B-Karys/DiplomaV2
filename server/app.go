package server

import (
	"DiplomaV2/internal/mailer"
	middleware2 "DiplomaV2/middleware"
	postHandlers "DiplomaV2/post/handlers"
	postModels "DiplomaV2/post/models"
	postRepositories "DiplomaV2/post/repository"
	postUseCases "DiplomaV2/post/usecase"
	userHandlers "DiplomaV2/user/handlers"
	userModels "DiplomaV2/user/models"
	userRepositories "DiplomaV2/user/repository"
	tokenRepositories "DiplomaV2/user/tokenRepository"
	userUseCases "DiplomaV2/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"

	"DiplomaV2/config"
	"DiplomaV2/database"
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
	appMailer := mailer.New("sandbox.smtp.mailtrap.io", 25, "f77e84f49aea4c", "ab28a0e2848b3b", "Test <no-reply@test.com>")

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
	// CORS middleware with configuration
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"}, // Specify your React frontend domain here
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
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

func (s *echoServer) initializePostHttpHandler() {
	postPostgresRepository := postRepositories.NewPostRepository(s.db)
	postUseCase := postUseCases.NewPostUseCase(postPostgresRepository)
	postHttpHandler := postHandlers.NewPostHttpHandler(postUseCase)

	postRouters := s.app.Group("/v2/posts")
	postRouters.POST("/", postHttpHandler.CreatePost, middleware2.LoginMiddleware)
}

func (s *echoServer) initializeUserHttpHandler() {
	userPostgresRepository := userRepositories.NewUserRepository(s.db)
	tokenPostgresRepository := tokenRepositories.NewTokenRepository(s.db)
	userUseCase := userUseCases.NewUserUseCase(userPostgresRepository, tokenPostgresRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(userUseCase, s.mailer)

	userRouters := s.app.Group("/v2/users")
	{
		userRouters.POST("/registration", userHttpHandler.Registration)
		userRouters.PUT("/", userHttpHandler.Activation)
		userRouters.POST("/login", userHttpHandler.Authentication)
		userRouters.GET("/check-auth", userHttpHandler.CheckAuth)
		userRouters.GET("/:id", userHttpHandler.GetUserInfoById, middleware2.LoginMiddleware)
		userRouters.PATCH("/update", userHttpHandler.UpdateUserInfo, middleware2.LoginMiddleware)
		userRouters.PATCH("/password", userHttpHandler.ChangePassword, middleware2.LoginMiddleware)
		userRouters.POST("/logout", userHttpHandler.Logout, middleware2.LoginMiddleware)
		userRouters.DELETE("/:id", userHttpHandler.DeleteUser, middleware2.LoginMiddleware)
		userRouters.POST("/forgot-password", userHttpHandler.ForgotPassword)
		userRouters.POST("/reset-password", userHttpHandler.ResetPassword)
	}
}

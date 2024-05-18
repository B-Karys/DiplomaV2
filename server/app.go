package server

import (
	"DiplomaV2/internal/mailer"
	middleware2 "DiplomaV2/middleware"
	postHandlers "DiplomaV2/post/handlers"
	postModels "DiplomaV2/post/models"
	postRepositories "DiplomaV2/post/repository"
	postUseCases "DiplomaV2/post/usecase"
	tokenRepositories "DiplomaV2/token/repository"
	userHandlers "DiplomaV2/user/handlers"
	userModels "DiplomaV2/user/models"
	userRepositories "DiplomaV2/user/repository"
	userUseCases "DiplomaV2/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

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
	s.app.Use(middleware.CORS())

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
		userRouters.POST("/", userHttpHandler.Registration)
		//userRouters.GET("/", userHttpHandler.GetUserInfo)
		//userRouters.PATCH("/", userHttpHandler.UpdateUserInfo, middleware2.LoginMiddleware)
		//userRouters.DELETE("/", userHttpHandler.DeleteUser, middleware2.LoginMiddleware)
		userRouters.PATCH("/", userHttpHandler.Activation, middleware2.LoginMiddleware)
		//userRouters.PATCH("/", userHttpHandler.ResetPassword)

	}
}

//func (s *echoServer) initializeTokenHttpHandler() {
//	tokenPostgresRepository := tokenRepositories.NewTokenRepository(s.db)
//	userPostgresRepository := userRepositories.NewUserRepository(s.db)
//	tokenUseCase := tokenUseCases.NewTokenUseCase(tokenPostgresRepository, userPostgresRepository)
//	tokenHttpHandler := tokenHandlers.NewTokenHttpHandler(tokenUseCase)
//
//	tokenRouters := s.app.Group("/v2/tokens")
//	{
//		tokenRouters.POST("/authentication", tokenHttpHandler.CreateAuthenticationToken)
//	}
//}

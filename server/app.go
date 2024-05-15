package server

import (
	middleware2 "DiplomaV2/middleware"
	postHandlers "DiplomaV2/post/handlers"
	postModels "DiplomaV2/post/models"
	postRepositories "DiplomaV2/post/repository"
	postUseCases "DiplomaV2/post/usecase"
	tokenHandlers "DiplomaV2/token/handlers"
	tokenModels "DiplomaV2/token/models"
	tokenRepositories "DiplomaV2/token/repository"
	tokenUseCases "DiplomaV2/token/usecase"
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
	app  *echo.Echo
	db   database.Database
	conf *config.Config
}

func NewEchoServer(conf *config.Config, db database.Database) Server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	return &echoServer{
		app:  echoApp,
		db:   db,
		conf: conf,
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

	// Initialize Post HTTP handler
	s.initializePostHttpHandler()
	s.initializeUserHttpHandler()
	s.initializeTokenHttpHandler()

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *echoServer) initializeMigrations() {
	s.db.GetDb().AutoMigrate(
		&userModels.User{},
		&postModels.Post{},
		&tokenModels.Token{},
	)
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
	userUseCase := userUseCases.NewUserUseCase(userPostgresRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(userUseCase)

	userRouters := s.app.Group("/v2/users")
	{
		userRouters.POST("/", userHttpHandler.Registration)
		//userRouters.GET("/", userHttpHandler.GetUserInfo)
		//userRouters.PATCH("/", userHttpHandler.UpdateUserInfo, middleware2.LoginMiddleware)
		//userRouters.DELETE("/", userHttpHandler.DeleteUser, middleware2.LoginMiddleware)
		//userRouters.PATCH("/", userHttpHandler.Activation, middleware2.LoginMiddleware)
		//userRouters.PATCH("/", userHttpHandler.ResetPassword)

	}
}

func (s *echoServer) initializeTokenHttpHandler() {
	tokenPostgresRepository := tokenRepositories.NewTokenRepository(s.db)
	userPostgresRepository := userRepositories.NewUserRepository(s.db)
	tokenUseCase := tokenUseCases.NewTokenUseCase(tokenPostgresRepository, userPostgresRepository)
	tokenHttpHandler := tokenHandlers.NewTokenHttpHandler(tokenUseCase)

	tokenRouters := s.app.Group("/v2/tokens")
	{
		tokenRouters.POST("/authentication", tokenHttpHandler.CreateAuthenticationToken)
	}
}

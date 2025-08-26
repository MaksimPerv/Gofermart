package app

import (
	"github.com/MaksimPerv/Gofermart/internal/config"
	"github.com/MaksimPerv/Gofermart/internal/controller/auth_handler"
	"github.com/MaksimPerv/Gofermart/internal/controller/user_handler"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	cfg         *config.Config
	router      *chi.Mux
	logger      *zap.Logger
	userService service.UserService
	authService service.AuthService
}

func New(cfg *config.Config, logger *zap.Logger, userService service.UserService, authService service.AuthService) *App {
	r := chi.NewRouter()
	return &App{cfg: cfg, router: r, logger: logger, userService: userService, authService: authService}

}

func (a *App) Router() http.Handler {
	return a.router
}

func (a *App) setupRoutes() {
	userHandler := user_handler.NewUserHandler(a.logger, a.userService)
	authHandler := auth_handler.NewAuthHandler(a.logger, a.authService)

	a.router.Post("/api/user/register", authHandler.Register)
	a.router.Post("/api/user/login", authHandler.Login)
	a.router.Get("/", userHandler.Get)
}

func (a *App) Run() error {
	a.setupRoutes()
	a.logger.Info("Server started", zap.String("port", a.cfg.RunAddress))
	return http.ListenAndServe(""+a.cfg.RunAddress, a.router)
}

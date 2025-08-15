package app

import (
	"github.com/MaksimPerv/Gofermart/internal/config"
	"github.com/MaksimPerv/Gofermart/internal/controller"
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
	userHandler := controller.NewUserHandler(a.logger, a.userService)
	authHandler := controller.NewAuthHandler(a.logger, a.authService)

	a.router.Post("/api/user/register", authHandler.Get)

	a.router.Get("/", userHandler.Get)
	a.router.Get("/123", authHandler.Get)
}

func (a *App) Run() error {
	a.setupRoutes()
	a.logger.Info("Server started", zap.String("port", a.cfg.RunAddress))
	return http.ListenAndServe(""+a.cfg.RunAddress, a.router)
}

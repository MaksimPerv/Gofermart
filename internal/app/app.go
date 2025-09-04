package app

import (
	"github.com/MaksimPerv/Gofermart/internal/config"
	"github.com/MaksimPerv/Gofermart/internal/controller/auth_handler"
	"github.com/MaksimPerv/Gofermart/internal/controller/order_handler"
	"github.com/MaksimPerv/Gofermart/internal/controller/user_handler"
	"github.com/MaksimPerv/Gofermart/internal/middleware"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	cfg          *config.Config
	router       *chi.Mux
	logger       *zap.Logger
	userService  service.UserService
	authService  service.AuthService
	orderService service.OrderService
}

func New(cfg *config.Config, logger *zap.Logger, userService service.UserService, authService service.AuthService, orderService service.OrderService) *App {
	r := chi.NewRouter()
	return &App{cfg: cfg, router: r, logger: logger, userService: userService, authService: authService, orderService: orderService}

}

func (a *App) Router() http.Handler {
	return a.router
}

func (a *App) setupRoutes() {
	userHandler := user_handler.NewUserHandler(a.logger, a.userService)
	authHandler := auth_handler.NewAuthHandler(a.logger, a.authService)
	orderHandler := order_handler.NewOrderHandler(a.logger, a.orderService)

	a.router.Group(func(r chi.Router) {
		r.Post("/api/user/register", authHandler.Register)
		r.Post("/api/user/login", authHandler.Login)
	})

	a.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(a.authService))

		r.Post("/api/user/orders", orderHandler.CreateOrder)
		r.Get("/api/user/orders", orderHandler.GetOrders)
		r.Get("/api/user/balance", userHandler.GetBalance)
		r.Post("/api/user/balance/withdraw", userHandler.WithdrawalRequest)

	})
}

func (a *App) Run() error {
	a.setupRoutes()
	a.logger.Info("Server started", zap.String("port", a.cfg.RunAddress))
	return http.ListenAndServe(""+a.cfg.RunAddress, a.router)
}

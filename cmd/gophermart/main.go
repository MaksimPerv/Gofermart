package main

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/app"
	"github.com/MaksimPerv/Gofermart/internal/config"
	"github.com/MaksimPerv/Gofermart/internal/database"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"github.com/MaksimPerv/Gofermart/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {

	cfg := config.Load()

	log := logger.New("debug")
	defer log.Sync()

	log.Info("Starting service with configuration",
		zap.String("run_address", cfg.RunAddress),
		zap.String("database_uri", cfg.DatabaseURI),
		zap.String("accrual_address", cfg.AccrualSystemAddress),
	)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err = db.Ping(ctx); err != nil {
		log.Fatal("DB ping failed", zap.Error(err))
	}

	if err = database.RunMigrations(db); err != nil {
		log.Fatal("Failed to roll out migrations", zap.Error(err))
	}

	userRepo := repository.NewPostgresUserRepository(db, log)
	orderRepo := repository.NewPostgresOrderRepository(db, log)
	orderService := service.NewOrderService(orderRepo, log)
	userService := service.NewUserService(userRepo, log)
	authService := service.NewAuthService(userRepo, log)
	app := app.New(cfg, log, userService, authService, orderService)

	if err = app.Run(); err != nil {
		log.Fatal("Server error", zap.Error(err))
	}

}

package service

import (
	"context"
	"github.com/MaksimPerv/Gofermart/internal/client/loyalty"
	"github.com/MaksimPerv/Gofermart/internal/repository"
	"go.uber.org/zap"
	"sync"
	"time"
)

type WorkerService struct {
	orderRepo     repository.OrderRepository
	loyaltyClient *loyalty.Client
	logger        *zap.Logger
	interval      time.Duration
	maxConcurrent int
}

func NewWorkerService(orderRepo repository.OrderRepository, loyaltyClient *loyalty.Client, logger *zap.Logger, interval time.Duration, n int) *WorkerService {
	return &WorkerService{
		orderRepo:     orderRepo,
		loyaltyClient: loyaltyClient,
		logger:        logger,
		interval:      interval,
		maxConcurrent: n,
	}
}

func (w *WorkerService) Start(ctx context.Context) {
	w.logger.Info("Start worker")
	ticket := time.NewTicker(w.interval)
	defer ticket.Stop()

	semaphore := make(chan struct{}, w.maxConcurrent)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopped")
			return
		case <-ticket.C:
			w.logger.Info("Start ticket")
			w.processOrders(ctx, semaphore)
			//	w.logger.Info("ЗАебись")

		}
	}
}

func (w *WorkerService) processOrders(ctx context.Context, semaphore chan struct{}) {
	orders, err := w.orderRepo.GetOrdersStatus(ctx)
	if err != nil {
		w.logger.Error("Failed to get orders to process", zap.Error(err))
		return
	}

	var wg sync.WaitGroup
	for _, orderNumber := range orders {
		//	w.logger.Info("Стартанул с ", zap.String("order", orderNumber))
		wg.Add(1)
		semaphore <- struct{}{}

		go func(order string) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()

			status, err := w.loyaltyClient.GetOrderStatus(ctx, order)
			if err != nil {
				w.logger.Error("Failed to get order status",
					zap.String("order", order),
					zap.Error(err))
				return
			}
			if status == nil {
				return
			}

			err = w.orderRepo.UpdateOrderStatus(ctx, status.Order, status.Status, status.Accrual)
			if err != nil {
				w.logger.Error("Failed to update order status",
					zap.String("order", status.Order),
					zap.Error(err))
				return
			}
			w.logger.Info("Order processed successfully",
				zap.String("order", status.Order),
				zap.String("status", status.Status))

		}(orderNumber)

	}
	wg.Wait()
	w.logger.Info("Закончил")
}

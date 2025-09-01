package order_handler

import (
	"errors"
	"github.com/MaksimPerv/Gofermart/internal/service"
	"github.com/MaksimPerv/Gofermart/internal/validation"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type OrderHandler struct {
	logger       *zap.Logger
	orderService service.OrderService
}

func NewOrderHandler(logger *zap.Logger, OrderService service.OrderService) *OrderHandler {
	return &OrderHandler{logger: logger, orderService: OrderService}
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	o.logger.Debug("CreateOrder request started",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))

	if r.Header.Get("Content-Type") != "text/plain" {
		o.logger.Warn("Invalid Content-Type",
			zap.String("content_type", r.Header.Get("Content-Type")))
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		o.logger.Error("User not authenticated in context")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		o.logger.Error("Failed to read request body", zap.Error(err))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))

	if orderNumber == "" {
		o.logger.Warn("Empty order number received")
		http.Error(w, "Order number is required", http.StatusBadRequest)
		return
	}

	if !validation.ValidateLuhn(orderNumber) {
		o.logger.Warn("Invalid order number format",
			zap.String("order_number", orderNumber))
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	o.logger.Info("Creating order in service",
		zap.String("order_number", orderNumber),
		zap.Int("user_id", userId))

	if err = o.orderService.CreateOrder(r.Context(), orderNumber, userId); err != nil {
		switch {
		case errors.Is(err, service.ErrOrderAlreadyUploaded):
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Order  already uploaded by this user"))
			return
		case errors.Is(err, service.ErrOrderBelongsToAnotherUser):
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	o.logger.Info("Order accepted for processing",
		zap.String("order_number", orderNumber),
		zap.Int("user_id", userId))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order accepted for processing"))

}
func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	o.logger.Debug("GetOrders request started",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))

	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		o.logger.Error("User not authenticated in context")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

}

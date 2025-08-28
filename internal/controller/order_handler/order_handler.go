package order_handler

import (
	"github.com/MaksimPerv/Gofermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type OrderHandler struct {
	logger       *zap.Logger
	orderService service.OrderService
}

func NewOrderHandler(logger *zap.Logger, OrderService service.OrderService) *OrderHandler {
	return &OrderHandler{logger: logger, orderService: OrderService}
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
	//	return
	//}
	//
	//userId, ok := r.Context().Value("userID").(int)
	//if !ok {
	//	http.Error(w, "User not authenticated", http.StatusUnauthorized)
	//	return
	//}
	//
	//defer r.Body.Close()
	//
	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	o.logger.Error("Failed to read request body", zap.Error(err))
	//	http.Error(w, "Failed to read request body", http.StatusBadRequest)
	//	return
	//}
	//
	//orderNumber := strings.TrimSpace(string(body))
	//
	//if orderNumber == "" {
	//	http.Error(w, "Order number is required", http.StatusBadRequest)
	//	return
	//}
	//
	//if !validation.ValidateLuhn(orderNumber) {
	//	http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
	//	return
	//}
	//
	//if err = o.orderService.CreateOrder(orderNumber); err != nil {
	//
	//}
}

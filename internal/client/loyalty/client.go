package loyalty

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type OrderStatus struct {
	Order   string          `json:"order"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetOrderStatus(ctx context.Context, orderNumber string) (*OrderStatus, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("error send request")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil, &ErrOrderNotRegistered{OrderNumber: orderNumber}
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")
		return nil, &ErrRateLimitExceeded{RetryAfter: retryAfter}
	case http.StatusInternalServerError:
		return nil, &ErrServerError{StatusCode: resp.StatusCode}
	case http.StatusOK:
		var status OrderStatus
		if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return nil, errors.New("JSON decode failed")
		}
		//log.Print(status.Status, status.Order, status.Accrual)
		return &status, nil
	default:
		return nil, err
	}

}

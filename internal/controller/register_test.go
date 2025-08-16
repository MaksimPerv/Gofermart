package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/MaksimPerv/Gofermart/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockAuthService) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockAuthService)
	handler := &AuthHandler{
		logger:      logger,
		authService: mockService,
	}

	t.Run("success", func(t *testing.T) {
		mockService.On("Register", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil).Once()
		mockService.On("GenerateToken", mock.AnythingOfType("*entity.User")).Return("test-token", nil).Once()

		user := map[string]string{
			"login":    "testuser",
			"password": "testpas",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		coolies := w.Result().Cookies()
		assert.Len(t, coolies, 1)
		assert.Equal(t, "token", coolies[0].Name)
		assert.Equal(t, "test-token", coolies[0].Value)

		assert.Contains(t, w.Body.String(), "User registered and authenticated")

		mockService.AssertExpectations(t)

	})
}

func TestRegister_ValidationErrors(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockAuthService)
	handler := &AuthHandler{
		logger:      logger,
		authService: mockService,
	}
	test := []struct {
		name     string
		input    map[string]string
		wantCode int
		wantBody string
	}{
		{
			name:     "Empty login",
			input:    map[string]string{"password": "testpass"},
			wantCode: http.StatusBadRequest,
			wantBody: "Login and password are required",
		},
		{
			name:     "Empty password",
			input:    map[string]string{"login": "testuser"},
			wantCode: http.StatusBadRequest,
			wantBody: "Login and password are required",
		},
		{
			name:     "Invalid content type",
			input:    map[string]string{"login": "testuser", "password": "testpass"},
			wantCode: http.StatusBadRequest,
			wantBody: "Content-Type must be application/json",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))

			if tt.name != "Invalid content type" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBody)

		})
	}
}

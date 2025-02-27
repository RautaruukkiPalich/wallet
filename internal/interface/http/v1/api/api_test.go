package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/internal/presenter"
	"wallet/internal/repository/wallet"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"wallet/internal/dto"
	"wallet/internal/interface/http/v1/api/mocks" // Путь к сгенерированным мокам
)

func TestPostOperation(t *testing.T) {
	router := mux.NewRouter()
	mockWallet := new(mocks.WalletPresenter) // Создание экземпляра нашего мока
	RegisterRouter(router, mockWallet)

	tests := []struct {
		name       string
		body       dto.PostOperationRequest
		statusCode int
		err        error
	}{
		{
			"Valid Request",
			dto.PostOperationRequest{
				WalletId:      "valid-uuid",
				OperationType: "credit",
				Amount:        100,
			},
			http.StatusOK,
			nil,
		},
		{
			"Invalid Wallet UUID",
			dto.PostOperationRequest{
				WalletId:      "invalid-uuid",
				OperationType: "credit",
				Amount:        100,
			},
			http.StatusBadRequest,
			presenter.ErrInvalidUUID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			if tt.name == "Valid Request" {
				// Настройка мока только для валидного запроса
				mockWallet.On("Transaction", mock.Anything, mock.Anything).Return(nil).Once()
			}
			if tt.name == "Invalid Wallet UUID" {
				// Настройка мока только для невалидного запроса
				mockWallet.On("Transaction", mock.Anything, mock.Anything).Return(presenter.ErrInvalidUUID).Once()
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			mockWallet.AssertExpectations(t) // Проверка, что все ожидания выполнены
		})
	}
}

func TestGetWalletAmount(t *testing.T) {
	router := mux.NewRouter()
	mockWallet := new(mocks.WalletPresenter)
	RegisterRouter(router, mockWallet)

	tcs := []struct {
		name       string
		uuid       string
		statusCode int
		expected   *dto.GetBalanceResponse
	}{
		{"Valid UUID", "valid-uuid", http.StatusOK, &dto.GetBalanceResponse{Amount: 100}},
		{"Invalid UUID", "invalid-uuid", http.StatusNotFound, nil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/wallets/"+tc.uuid, nil)
			w := httptest.NewRecorder()

			if tc.name == "Valid UUID" {
				mockWallet.On("GetBalance", mock.Anything, tc.uuid).Return(&dto.GetBalanceResponse{Amount: 100}, nil).Once()
			} else {
				mockWallet.On("GetBalance", mock.Anything, tc.uuid).Return(nil, wallet.ErrWalletNotFound).Once()
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.statusCode, w.Code)
			if tc.expected != nil {
				var response dto.GetBalanceResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expected.Amount, response.Amount)
			}
			mockWallet.AssertExpectations(t)
		})
	}
}

func TestCreateWallet(t *testing.T) {
	router := mux.NewRouter()
	mockWallet := new(mocks.WalletPresenter)
	RegisterRouter(router, mockWallet)

	req := httptest.NewRequest(http.MethodPost, "/wallet/create", nil)
	w := httptest.NewRecorder()

	mockWallet.On("NewWallet", mock.Anything).Return(&dto.WalletResponse{UUID: "new-uuid", Amount: 0}, nil).Once()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := new(dto.WalletResponse)
	if err := json.NewDecoder(w.Body).Decode(response); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "new-uuid", response.UUID)
	assert.Equal(t, int64(0), response.Amount)

	mockWallet.AssertExpectations(t)
}

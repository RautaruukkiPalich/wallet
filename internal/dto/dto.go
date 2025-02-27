package dto

type GetBalanceResponse struct {
	Amount int64 `json:"amount"`
}

type WalletResponse struct {
	UUID   string `json:"uuid"`
	Amount int64  `json:"amount"`
}

type PostOperationRequest struct {
	WalletId      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

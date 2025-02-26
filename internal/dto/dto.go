package dto

type GetBalanceResponse struct {
	Amount int64 `json:"amount"`
}

type WalletResponse struct {
	UUID   string `json:"uuid"`
	Amount int64  `json:"amount"`
}

type PostOperationRequest struct {
	WalletUUID    string `json:"walletUUID"` // так как используем UUID, поле решил тоже назвать UUID, а не ID
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

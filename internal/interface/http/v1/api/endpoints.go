package api

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"wallet/internal/dto"
	"wallet/internal/interface/response"
)

// @Summary		PostOperation
// @Description	add operation using walletUUID, operationType and amount
// @Tags			wallets
// @Accept			json
// @Produce		json
// @Param			input	body		dto.PostOperationRequest	true	"request"
// @Success		200		{object}	nil
// @Failure		400,404	{object}	dto.ErrorResponse
// @Success		500		{object}	dto.ErrorResponse
// @Success		default	{object}	dto.ErrorResponse
// @Router			/ [post]
func (rt *Router) postOperation(w http.ResponseWriter, r *http.Request) {
	var req dto.PostOperationRequest

	if err := getFromBody(r, &req); err != nil {
		response.
			Resp().
			WithCode(http.StatusBadRequest).
			WithError(ErrInvalidFormData).
			Build().
			Write(w)
		return
	}

	if err := rt.wallet.Transaction(context.TODO(), &req); err != nil {
		response.
			Resp().
			HandleError(err).
			Build().
			Write(w)
		return
	}

	response.Resp().WithCode(http.StatusOK).Build().Write(w)
}

// @Summary		GetAmount
// @Description	get amount by wallets uuid
// @Tags			wallets
// @Accept			json
// @Produce		json
// @Param			uuid	path		string	false	"wallet_uuid"
// @Success		200		{object}	dto.GetBalanceResponse
// @Failure		400,404	{object}	dto.ErrorResponse
// @Success		500		{object}	dto.ErrorResponse
// @Success		default	{object}	dto.ErrorResponse
// @Router			/{uuid} [get]
func (rt *Router) getWalletAmount(w http.ResponseWriter, r *http.Request) {
	const uuid = "uuid"
	walletUUID := mux.Vars(r)[uuid]

	if walletUUID == "" {
		response.
			Resp().
			WithCode(http.StatusBadRequest).
			WithError(ErrEmptyWalletUUID).
			Build().
			Write(w)
		return
	}

	balance, err := rt.wallet.GetBalance(context.TODO(), walletUUID)
	if err != nil {
		response.Resp().HandleError(err).Build().Write(w)
		return
	}

	response.Resp().WithCode(http.StatusOK).WithPayload(balance).Build().Write(w)
}

// @Summary		CreateWallet
// @Description	create new wallet, returning uuid and amount
// @Tags			wallets
// @Accept			json
// @Produce		json
// @Success		200		{object}	dto.WalletResponse
// @Failure		400,404	{object}	dto.ErrorResponse
// @Success		500		{object}	dto.ErrorResponse
// @Success		default	{object}	dto.ErrorResponse
// @Router			/create [post]
func (rt *Router) createWallet(w http.ResponseWriter, r *http.Request) {
	wallet, err := rt.wallet.NewWallet(context.TODO())
	if err != nil {
		response.Resp().HandleError(err).Build().Write(w)
		return
	}

	response.Resp().WithCode(http.StatusOK).WithPayload(wallet).Build().Write(w)
}

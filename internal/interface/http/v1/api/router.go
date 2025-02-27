package api

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"wallet/internal/dto"
)

//go:generate mockery --name walletPresenter --structname=WalletPresenter
type walletPresenter interface {
	Transaction(context.Context, *dto.PostOperationRequest) error
	GetBalance(context.Context, string) (*dto.GetBalanceResponse, error)
	NewWallet(ctx context.Context) (*dto.WalletResponse, error)
}

type Router struct {
	wallet walletPresenter
	router *mux.Router
}

const (
	postOperationPath   = "/wallet"
	createWalletPath    = "/wallet/create"
	getWalletAmountPath = "/wallets/{uuid}"
)

func RegisterRouter(
	router *mux.Router,
	wallet walletPresenter,
) *Router {
	rt := &Router{
		router: router,
		wallet: wallet,
	}

	rt.router.HandleFunc(postOperationPath, rt.postOperation).Methods(http.MethodPost)
	rt.router.HandleFunc(createWalletPath, rt.createWallet).Methods(http.MethodPost)
	rt.router.HandleFunc(getWalletAmountPath, rt.getWalletAmount).Methods(http.MethodGet)

	return rt
}

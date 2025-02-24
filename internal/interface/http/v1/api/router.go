package api

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"wallet/internal/dto"
)

type walletPresenter interface {
	Transaction(context.Context, *dto.PostOperationRequest) error
	GetBalance(context.Context, string) (*dto.GetBalanceResponse, error)
}

type Router struct {
	wallet walletPresenter
	router *mux.Router
}

const (
	postOperationPath   = "/"
	getWalletAmountPath = "/{uuid}"
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
	rt.router.HandleFunc(getWalletAmountPath, rt.getWalletAmount).Methods(http.MethodGet)

	return rt
}

package app

import (
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wallet/config"
	"wallet/internal/interface/http/v1/api"
	"wallet/internal/presenter"
	"wallet/internal/services"
	"wallet/internal/utils/httpserver"
	"wallet/internal/utils/metrics"
	"wallet/internal/utils/pprof"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "wallet/docs" // docs generated by Swag CLI, you have to import it.
)

const (
	pathToAPIWallets = "/api/v1/wallets"
)

func Run(cfg *config.Config) {

	profilerServer := pprof.NewPProfServer(cfg.PProf.Convert())
	defer func() {
		if err := profilerServer.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
	go profilerServer.Run()

	metricsServer := metrics.NewMetricsServer(cfg.Metrics.Convert())
	defer func() {
		if err := metricsServer.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
	go metricsServer.Run()

	walletService := services.NewWalletService()

	walletPresenter := presenter.NewPresenter(walletService)

	router := mux.NewRouter()
	router.Use(
		metrics.MW,
	)

	walletRouter := router.PathPrefix(pathToAPIWallets).Subrouter()
	walletRouter.PathPrefix("/swagger/").HandlerFunc(httpSwagger.WrapHandler)
	api.RegisterRouter(walletRouter, walletPresenter)

	server := httpserver.NewHTTPServer(cfg.HTTPServer.Convert(), router)
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
	go server.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	select {
	case res := <-stop:
		log.Println("syscall stop", res.String())
	case err := <-server.Notify():
		log.Println("http server notify: ", err)
	case err := <-metricsServer.Notify():
		log.Println("metrics notify: ", err)
	case err := <-profilerServer.Notify():
		log.Println("profiler notify: ", err)
	}

	log.Println("service exit")
}

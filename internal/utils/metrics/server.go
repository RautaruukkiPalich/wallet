package metrics

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Server struct {
	srv    *http.Server
	router *mux.Router
	notify chan error
}

type Config struct {
	Addr string
}

func NewMetricsServer(cfg Config) *Server {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())

	return &Server{
		&http.Server{
			Addr:    cfg.Addr,
			Handler: r,
		},
		r,
		make(chan error, 1),
	}
}

func (s *Server) Run() {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				log.Printf("metrics - server error : %s", v.Error())
			default:
				log.Printf("metrics - server fail : %v", v)
			}
		}
	}()

	s.notify <- s.srv.ListenAndServe()
	close(s.notify)
}

func (s *Server) Shutdown() error {
	return s.srv.Shutdown(context.Background())
}

func (s *Server) Notify() chan error {
	return s.notify
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

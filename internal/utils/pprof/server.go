package pprof

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/pprof"
)

type Server struct {
	srv    *http.Server
	router *mux.Router
	notify chan error
}

type Config struct {
	Addr string
}

func NewPProfServer(cfg Config) *Server {
	r := mux.NewRouter()

	sr := r.PathPrefix("/debug/pprof").Subrouter()

	sr.HandleFunc("/", pprof.Index)
	sr.HandleFunc("/cmdline", pprof.Cmdline)
	sr.HandleFunc("/profile", pprof.Profile)
	sr.HandleFunc("/symbol", pprof.Symbol)
	sr.HandleFunc("/trace", pprof.Trace)
	sr.HandleFunc("/{profile}", pprof.Index)

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
				log.Printf("pprof - server error : %s\n", v.Error())

			default:
				log.Printf("pprof - server error : %v\n", v)
			}
		}
	}()

	s.notify <- s.srv.ListenAndServe()
	close(s.notify)
}

func (s *Server) Shutdown() error {
	return s.srv.Shutdown(context.TODO())
}

func (s *Server) Notify() chan error {
	return s.notify
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

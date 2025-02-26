package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Server struct {
	srv       http.Server
	router    http.Handler
	notify    chan error
	isRunning atomic.Bool
}

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewHTTPServer(cfg ServerConfig, router http.Handler) *Server {
	return &Server{
		srv: http.Server{
			Addr:         cfg.Addr,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		router: router,
		notify: make(chan error),
	}
}

func (s *Server) Run() {
	if s.isRunning.Swap(true) {
		fmt.Println(ErrDuplicateRun)
		return
	}

	s.notify <- s.srv.ListenAndServe()
	close(s.notify)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Notify() chan error {
	return s.notify
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

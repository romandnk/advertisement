package http

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(host string, port int, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	srv := &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return &Server{srv: srv}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

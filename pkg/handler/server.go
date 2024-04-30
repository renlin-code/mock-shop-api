package handler

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httServer.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httServer.Shutdown(ctx)
}

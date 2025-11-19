package serser

import (
	"context"
	"net/http"
	"time"

	"github.com/kcthack-auth/internal/config"
)

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler, config config.Config) *Server {
	return &Server{server: &http.Server{
		Addr:           ":" + config.HTTP.Port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

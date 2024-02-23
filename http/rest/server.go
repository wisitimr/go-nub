package rest

import (
	"context"
	"findigitalservice/config"
	"findigitalservice/pkg/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg    config.HTTPServer
	router *chi.Mux
	logger *logrus.Logger
}

func NewServer(ctx context.Context, cfg config.Configuration) (*Server, error) {
	mongodb, err := db.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Can't establish database connection")
	}
	logger := NewLogger()
	s := Server{
		logger: logger,
		cfg:    cfg.HTTPServer,
		router: Register(mongodb, logger),
	}
	return &s, nil
}

func (s Server) Start(ctx context.Context) error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		IdleTimeout:  s.cfg.IdleTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("Server started on port %d", s.cfg.Port)
		serverErrors <- server.ListenAndServe()
	}(&wg)

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting rest api server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received stop signal")
		// asking listener to shutdown
		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait()
		s.logger.Info("server was shut down gracefully")
	}
	return nil
}

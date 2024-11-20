package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-app/internal/application"
	"wallet-app/internal/storage"
	"wallet-app/pkg/logger"
)

func Run(cfg *Config, log *logger.Log) error {
	db, err := storage.NewPostgresDB(cfg.DatabaseUrl)
	if err != nil {
		log.Error("error connecting to database", err)
		os.Exit(1)
	}
	defer db.Close()

	httpServer := application.NewHttpServer(db)
	var serverAddress = cfg.ListenAndPort()
	server := &http.Server{
		Addr: serverAddress,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpServer.ServeHTTP(w, r)
		}),
		IdleTimeout:  cfg.ServerIdleTimeout,
		ReadTimeout:  cfg.ServerTimeout,
		WriteTimeout: cfg.ServerTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info(fmt.Sprintf("Starting server at %s", serverAddress))
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start http server", err)
		}
	}()
	log.Info("Server started and ready to accept connections")

	<-done
	log.Info("Server stopped")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	log.Info("Server shutdown properly")
	return nil
}

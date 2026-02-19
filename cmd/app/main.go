package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toggle-interviews/tynerj-interview-test/internal/handler"
	"github.com/toggle-interviews/tynerj-interview-test/internal/repo"
	"go.uber.org/zap"
)

// Data source for the staff records
const csvStaffFileName = "staff.csv"
const addr = ":8080"

func main() {
	if err := run(); err != nil {
		// Last-resort logging
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Load data at start up and fail if it is non-existent or invalid
	staffRepo, err := repo.NewCSVStaffRepo(csvStaffFileName)
	if err != nil {
		logger.Fatal("failed to load repository", zap.String("file", csvStaffFileName), zap.Error(err))
	}
	logger.Info("staff repository has been successfully loaded")

	r := gin.Default()
	// Disable trusted proxies to avoid trusting forwarded headers;
	_ = r.SetTrustedProxies(nil)

	healthHandler := handler.NewHealthHandler()
	staffHandler := handler.NewStaffHandler(staffRepo)

	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/staff/:id", staffHandler.GetStaffMemberByID)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start and run server in the background
	go func() {
		logger.Info("listening", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Can't return errors from goroutine so they have to be logged log it.
			logger.Error("server failed", zap.Error(err))
		}
	}()

	// Block until a shut down signal is received from os (i.e., CTRL+C, or kill <pid>)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")

	// Create a timelimit for how long willing to wait for things to shut down before force exiting
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Allows for a graceful shutdown; stops new connections, while allowing for logs and in progress requests to finish
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
		return fmt.Errorf("shutdown: %w", err)
	}

	logger.Info("server stopped cleanly")
	return nil
}

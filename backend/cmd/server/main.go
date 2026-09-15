package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IshaanShivKr/urlvio/internal/auth"
	"github.com/IshaanShivKr/urlvio/internal/config"
	"github.com/IshaanShivKr/urlvio/internal/database"
	"github.com/IshaanShivKr/urlvio/internal/handler"
	"github.com/IshaanShivKr/urlvio/internal/middleware"
	"github.com/IshaanShivKr/urlvio/internal/repository"
	"github.com/IshaanShivKr/urlvio/internal/routes"
	"github.com/IshaanShivKr/urlvio/internal/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	dbConnectTimeout  = 5 * time.Second
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	clerk.SetKey(cfg.ClerkSecretKey)

	dbCtx, dbCancel := context.WithTimeout(
		context.Background(),
		dbConnectTimeout,
	)
	defer dbCancel()

	db, err := database.NewPostgresPool(dbCtx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	slog.Info("connected to database")

	gin.SetMode(cfg.GinMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://urlvio.vercel.app/",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.SecurityHeaders())

	healthHandler := handler.NewHealthHandler(db)

	repo := repository.NewPostgresRepository(db)
	linkService := service.NewLinkService(repo)
	linkHandler := handler.NewLinkHandler(linkService, cfg.BaseURL)

	routes.Register(router, healthHandler, linkHandler, auth.Middleware())

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	slog.Info("starting server", "addr", server.Addr)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	return waitForShutdown(server, serverErr)
}

func waitForShutdown(server *http.Server, serverErr <-chan error) error {
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil

	case <-signalCtx.Done():
		return shutdown(server)
	}
}

func shutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	slog.Info("server stopped")
	return nil
}

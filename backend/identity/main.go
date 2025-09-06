package main

import (
	"context"
	"errors"
	"fmt"
	"krypton/identity/cmd"
	"krypton/identity/config"
	"krypton/identity/internal/adapters"
	"krypton/identity/internal/application/services"
	"krypton/identity/internal/delivery/http/error_mapping"
	httpserver "krypton/identity/internal/delivery/http/router"
	"krypton/pkg/auth"
	"krypton/pkg/gin/error_handler"
	"krypton/pkg/logger"
	"krypton/pkg/minio"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	rootCmd := cmd.NewRootCmd(runServer)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServer() error {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Init("development")
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(logger.Logger)

	appLogger := logger.Logger
	defer func(appLogger *zap.Logger) {
		_ = appLogger.Sync()
	}(appLogger)

	ctx := context.Background()
	// --- Database Connections ---
	appLogger.Info("connecting to Postgresql...")
	pgConn, err := pgx.Connect(ctx, cfg.Database.Postgresql.URI)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer func(pgConn *pgx.Conn, ctx context.Context) {
		_ = pgConn.Close(ctx)
	}(pgConn, ctx)
	appLogger.Info("Successfully connected to PostgreSQL")

	// --- Adapters ---
	adapter, err := adapters.NewAdapter(pgConn, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create adapters", zap.Error(err))
	}

	// --- Application Services ---
	tokenManager := auth.NewTokenManager(cfg.Auth.JWTSecretKey)
	service := services.NewService(adapter.UserRepo, tokenManager, appLogger)

	errCh := make(chan error, 1) // TODO Add 2 for grpc after add

	// Setup Gin
	go func() {
		errCh <- runHTTPServer(cfg.General.HTTPPort, service.UserService, nil, cfg, tokenManager, appLogger)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("a server failed: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func runHTTPServer(port string, userService services.UserService, minioService *minio.Service, cfg *config.Config, tokenManager *auth.TokenManager, appLogger *zap.Logger) error {
	appLogger.Info("starting HTTP (Gin) server...", zap.String("port", port))
	gin.SetMode(gin.DebugMode)
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	domainErrorMappings := error_mapping.GetDomainErrorMappings()
	engine.Use(error_handler.New(domainErrorMappings, appLogger))

	httpserver.Setup(engine, userService, minioService, cfg, tokenManager, appLogger)

	appLogger.Info("HTTP (Gin) Server is running on", zap.String("port", port))
	if err := engine.Run(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("gin server failed: %w", err)
	}
	return nil
}

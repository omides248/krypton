package main

import (
	"context"
	"errors"
	"fmt"
	"krypton/pkg/auth"
	"krypton/pkg/gin/error_handler"
	"krypton/pkg/logger"
	"krypton/wallet/cmd"
	"krypton/wallet/config"
	"krypton/wallet/internal/adapters/storage/postgresql"
	"krypton/wallet/internal/application/services"
	"krypton/wallet/internal/delivery/http/error_mapping"
	httpserver "krypton/wallet/internal/delivery/http/router"
	"log"
	"net/http"
	"os"
	"time"

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
	appLogger := logger.Logger
	defer func(appLogger *zap.Logger) {
		_ = appLogger.Sync()
	}(appLogger)

	ctx := context.Background() // TODO timeout
	// --- Database Connection ---
	appLogger.Info("connecting to Postgresql...")
	pgConn, err := pgx.Connect(ctx, cfg.Database.Postgresql.URI)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer func(pgConn *pgx.Conn, ctx context.Context) {
		_ = pgConn.Close(ctx)
	}(pgConn, ctx)
	appLogger.Info("Successfully connected to PostgreSQL")

	// --- Repositories ---
	dbTimeout := 10 * time.Second
	assetRepo, _ := postgresql.NewAssetRepository(pgConn, dbTimeout, appLogger)
	accountRepo := postgresql.NewAccountRepository(pgConn, dbTimeout, appLogger)

	// --- Application Services ---
	tokenManager := auth.NewTokenManager(cfg.Auth.JWTSecretKey)
	walletService := services.NewWalletService(accountRepo, assetRepo, appLogger)

	errCh := make(chan error, 1)

	// --- Setup Gin ---
	go func() {
		errCh <- runHTTPServer(cfg.General.HTTPPort, walletService, tokenManager, appLogger)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("a server failed: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func runHTTPServer(port string, walletService services.WalletService, tokenManager *auth.TokenManager, appLogger *zap.Logger) error {
	appLogger.Info("starting HTTP (Gin) server...", zap.String("port", port))
	gin.SetMode(gin.DebugMode)
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	// Use custom error handler middleware
	domainErrorMappings := error_mapping.GetDomainErrorMappings()
	engine.Use(error_handler.New(domainErrorMappings, appLogger))

	// Setup routes
	httpserver.Setup(engine, walletService, tokenManager, appLogger)

	appLogger.Info("HTTP (Gin) Server is running on", zap.String("port", port))
	if err := engine.Run(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("gin server failed: %w", err)
	}
	return nil
}

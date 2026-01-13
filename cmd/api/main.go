package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UDL-TF/UnitedAPI/internal/config"
	appctx "github.com/UDL-TF/UnitedAPI/internal/context"
	"github.com/UDL-TF/UnitedAPI/internal/database"
	"github.com/UDL-TF/UnitedAPI/internal/logger"
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/UDL-TF/UnitedAPI/internal/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Server.Environment); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create application context for dependency injection
	ctx := context.Background()
	appContext := appctx.NewAppContext(ctx)

	db, err := database.InitPostgres(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	appContext.SetDB(db)

	// Create Gin engine
	engine := gin.New()

	// Apply global middlewares
	engine.Use(middleware.Recovery())                   // Recover from panics
	engine.Use(middleware.Logger())                     // Log all requests
	engine.Use(middleware.CORS())                       // Enable CORS
	engine.Use(middleware.InjectAppContext(appContext)) // Inject app context

	// Setup routes with versioning
	router.SetupRouter(engine)

	// Create HTTP server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		logger.Log.Info("Server starting",
			zap.String("port", cfg.Server.Port),
			zap.String("environment", cfg.Server.Environment))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	if db := appContext.GetDB(); db != nil {
		database.Close(db)
	}

	logger.Log.Info("Server exited gracefully")
}

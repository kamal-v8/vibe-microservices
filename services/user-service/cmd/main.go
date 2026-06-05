// User Service — entry point for the Pulse social platform's user management API.
//
// This service provides CRUD operations for user profiles and is designed to run
// as a containerized microservice behind an API gateway or load balancer.
//
// Architecture: cmd/main.go → handler → repository → PostgreSQL
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver — imported for side-effect registration

	"github.com/pulse-social/user-service/internal/config"
	"github.com/pulse-social/user-service/internal/handler"
	"github.com/pulse-social/user-service/internal/middleware"
	"github.com/pulse-social/user-service/internal/repository"
)

func main() {
	// ── Load configuration ──────────────────────────────────────────────
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal(`{"level":"fatal","msg":"DATABASE_URL environment variable is required","service":"user-service"}`)
	}

	logJSON("info", "starting user service", fmt.Sprintf("port=%s", cfg.Port))

	// ── Connect to PostgreSQL with retry ────────────────────────────────
	// Retries handle the common case where the database container is still
	// starting when this service comes up (e.g. in docker-compose).
	db, err := connectWithRetry(cfg.DatabaseURL, 5, 3*time.Second)
	if err != nil {
		log.Fatalf(`{"level":"fatal","msg":"failed to connect to database","error":"%v","service":"user-service"}`, err)
	}
	defer db.Close()

	logJSON("info", "connected to database", "")

	// ── Initialize layers ───────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// ── Set up Gin router ───────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Recovery middleware to catch panics and return 500 instead of crashing
	router.Use(gin.Recovery())
	// Structured JSON request logging
	router.Use(middleware.JSONLogger())

	// ── Register routes ─────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", userHandler.HealthCheck)

		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users", userHandler.ListUsers)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Metrics endpoint lives outside the /api/v1 prefix (Prometheus convention)
	router.GET("/metrics", userHandler.Metrics)

	// ── Start HTTP server ───────────────────────────────────────────────
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start serving in a goroutine so we can listen for shutdown signals
	go func() {
		logJSON("info", fmt.Sprintf("listening on :%s", cfg.Port), "")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(`{"level":"fatal","msg":"server failed","error":"%v","service":"user-service"}`, err)
		}
	}()

	// ── Graceful shutdown ───────────────────────────────────────────────
	// Wait for SIGTERM (Kubernetes pod termination) or SIGINT (Ctrl+C).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logJSON("info", fmt.Sprintf("received signal %v, shutting down gracefully", sig), "")

	// Give in-flight requests up to 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf(`{"level":"fatal","msg":"forced shutdown","error":"%v","service":"user-service"}`, err)
	}

	logJSON("info", "server stopped cleanly", "")
}

// connectWithRetry attempts to open and ping a PostgreSQL connection, retrying
// on failure with a fixed backoff. This handles the startup-ordering race
// condition that's common in containerized environments.
func connectWithRetry(dsn string, maxRetries int, backoff time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			logJSON("warn", fmt.Sprintf("database connection attempt %d/%d failed (open)", i+1, maxRetries), err.Error())
			time.Sleep(backoff)
			continue
		}

		// sql.Open only validates the DSN; Ping actually connects
		if err = db.Ping(); err != nil {
			logJSON("warn", fmt.Sprintf("database connection attempt %d/%d failed (ping)", i+1, maxRetries), err.Error())
			db.Close()
			time.Sleep(backoff)
			continue
		}

		// Connection successful — configure the pool
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		return db, nil
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

// logJSON emits a structured JSON log line to stdout.
func logJSON(level, msg, detail string) {
	ts := time.Now().UTC().Format(time.RFC3339)
	if detail != "" {
		fmt.Printf(`{"level":"%s","msg":"%s","detail":"%s","timestamp":"%s","service":"user-service"}`+"\n", level, msg, detail, ts)
	} else {
		fmt.Printf(`{"level":"%s","msg":"%s","timestamp":"%s","service":"user-service"}`+"\n", level, msg, ts)
	}
}

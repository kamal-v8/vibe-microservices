// Package handler implements the HTTP handlers (controllers) for the user API.
// Each handler method maps to a single REST endpoint and follows a consistent
// JSON envelope: {"data": ..., "error": ...}.
package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pulse-social/user-service/internal/model"
	"github.com/pulse-social/user-service/internal/repository"
)

// UserHandler holds dependencies needed by the HTTP handlers.
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler creates a handler wired to the given repository.
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// ---------- Response helpers ----------

// jsonResponse is the consistent envelope for all API responses.
type jsonResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func successResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, jsonResponse{Data: data, Error: nil})
}

func errorResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, jsonResponse{Data: nil, Error: msg})
}

// ---------- Handlers ----------

// CreateUser handles POST /api/v1/users.
// Validates the request body, inserts a new user, and returns 201 Created.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		// Check for PostgreSQL unique-violation (code 23505) to return a
		// user-friendly error instead of a generic 500.
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			errorResponse(c, http.StatusConflict, "username or email already exists")
			return
		}
		log.Printf(`{"level":"error","msg":"failed to create user","error":"%v","service":"user-service"}`, err)
		errorResponse(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	successResponse(c, http.StatusCreated, user)
}

// GetUser handles GET /api/v1/users/:id.
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format before hitting the database
	if _, err := uuid.Parse(id); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user ID format")
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		log.Printf(`{"level":"error","msg":"failed to get user","error":"%v","service":"user-service"}`, err)
		errorResponse(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	successResponse(c, http.StatusOK, user)
}

// ListUsers handles GET /api/v1/users with optional ?limit= and ?offset= params.
// Defaults: limit=20 (max 100), offset=0.
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Clamp limit to a reasonable maximum to prevent abuse
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		log.Printf(`{"level":"error","msg":"failed to list users","error":"%v","service":"user-service"}`, err)
		errorResponse(c, http.StatusInternalServerError, "failed to list users")
		return
	}

	successResponse(c, http.StatusOK, users)
}

// UpdateUser handles PUT /api/v1/users/:id (partial update).
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			errorResponse(c, http.StatusConflict, "username or email already exists")
			return
		}
		log.Printf(`{"level":"error","msg":"failed to update user","error":"%v","service":"user-service"}`, err)
		errorResponse(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	successResponse(c, http.StatusOK, user)
}

// DeleteUser handles DELETE /api/v1/users/:id.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user ID format")
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		log.Printf(`{"level":"error","msg":"failed to delete user","error":"%v","service":"user-service"}`, err)
		errorResponse(c, http.StatusInternalServerError, "failed to delete user")
		return
	}

	// 204 No Content — no response body
	c.Status(http.StatusNoContent)
}

// HealthCheck handles GET /api/v1/health.
// Pings the database to verify full-stack connectivity, not just the HTTP layer.
func (h *UserHandler) HealthCheck(c *gin.Context) {
	status := "ok"
	httpStatus := http.StatusOK

	if err := h.repo.Ping(c.Request.Context()); err != nil {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
		log.Printf(`{"level":"warn","msg":"health check database ping failed","error":"%v","service":"user-service"}`, err)
	}

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"service":   "user-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "v3-canary",
	})
}

// Metrics handles GET /metrics — a stub for future Prometheus integration.
func (h *UserHandler) Metrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "metrics endpoint - integrate with Prometheus",
	})
}

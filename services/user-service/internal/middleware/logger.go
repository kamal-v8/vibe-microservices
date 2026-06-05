// Package middleware provides Gin middleware functions for the user service.
package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONLogger returns a Gin middleware that logs every HTTP request in structured
// JSON format. This is intentionally a simple Printf-based approach to avoid
// pulling in a third-party logging library — in production you'd use zerolog
// or zap, but this keeps the dependency surface small for the exercise.
func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate latency after the request has been handled
		latency := time.Since(start)

		// Structured JSON log line — one line per request, easy to parse with
		// tools like jq, Loki, or CloudWatch Logs Insights.
		fmt.Printf(`{"level":"info","msg":"request completed","method":"%s","path":"%s","status":%d,"latency_ms":%d,"client_ip":"%s","timestamp":"%s","service":"user-service"}`+"\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency.Milliseconds(),
			c.ClientIP(),
			time.Now().UTC().Format(time.RFC3339),
		)
	}
}

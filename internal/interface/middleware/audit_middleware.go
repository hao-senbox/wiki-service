package middleware

import (
	"time"

	"wiki-service/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// AuditMiddleware logs all requests for audit purposes
type AuditMiddleware struct {
	logger *logger.Logger
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(logger *logger.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		logger: logger,
	}
}

// Log logs the request for audit purposes
func (m *AuditMiddleware) Log() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Read request payload
		reqBody := c.Body()

		// Process request
		err := c.Next()

		duration := time.Since(start)

		// Response payload
		resBody := c.Response().Body()
		statusCode := c.Response().StatusCode()
		ip := c.IP()

		// Log HTTP (bạn giữ nguyên)
		m.logger.HTTP(c.Method(), c.Path(), ip, statusCode, duration, nil)

		// Check if should audit
		if shouldAudit(c.Path(), c.Method()) {
			m.logger.Audit(
				c.Method(),      // method
				c.Path(),        // endpoint
				string(reqBody), // payload
				string(resBody), // response
				ip,              // ip
				statusCode,      // status
			)
		}

		return err
	}
}

// shouldAudit determines if a request should be audited
func shouldAudit(path, method string) bool {
	// Audit important endpoints
	auditPaths := []string{
		"/api/v1/wikis/template",
	}

	// Don't audit health checks and static files
	if path == "/health" || path == "/metrics" {
		return false
	}

	// Audit POST, PUT, DELETE operations
	if method == "POST" || method == "PUT" || method == "DELETE" {
		return true
	}

	// Check if path matches audit paths
	for _, auditPath := range auditPaths {
		if contains(path, auditPath) {
			return true
		}
	}

	return false
}

// determineAction determines the action based on path and method
func determineAction(path, method string) string {
	// Map common paths to actions
	pathActions := map[string]string{
		"/api/v1/auth/register": "user.register",
		"/api/v1/auth/login":    "user.login",
		"/api/v1/auth/logout":   "user.logout",
		"/api/v1/auth/profile":  "user.view_profile",
	}

	if action, ok := pathActions[path]; ok {
		return action
	}

	// Default action based on method
	switch method {
	case "GET":
		return "resource.view"
	case "POST":
		return "resource.create"
	case "PUT", "PATCH":
		return "resource.update"
	case "DELETE":
		return "resource.delete"
	default:
		return "resource.access"
	}
}

// determineResource determines the resource type from path
func determineResource(path string) string {
	if contains(path, "/auth") {
		return "auth"
	}
	if contains(path, "/users") {
		return "user"
	}
	if contains(path, "/profile") {
		return "profile"
	}
	return "unknown"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

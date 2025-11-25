package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"wiki-service/pkg/logger"
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

		// Process request first
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user ID from context if authenticated
		var userID *int64
		if uid := c.Locals("user_id"); uid != nil {
			if id, ok := uid.(int64); ok {
				userID = &id
			}
		}

		// Determine action and resource
		action := determineAction(c.Path(), c.Method())
		resource := determineResource(c.Path())

		// Get client info
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")
		statusCode := c.Response().StatusCode()

		// Log HTTP request
		m.logger.HTTP(c.Method(), c.Path(), ipAddress, statusCode, duration, userID)

		// Log audit entry for important actions
		if shouldAudit(c.Path(), c.Method()) {
			errorMsg := ""
			if statusCode >= 400 {
				errorMsg = "request failed"
			}

			m.logger.Audit(
				userID,
				action,
				resource,
				c.Method(),
				c.Path(),
				ipAddress,
				userAgent,
				statusCode,
				errorMsg,
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

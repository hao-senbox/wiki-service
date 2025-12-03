package libs_helper

import (
	"wiki-service/logger"

	"github.com/gofiber/fiber/v2"
)

const (
	ErrInvalidOperation = "ERR_INVALID_OPERATION"
	ErrInvalidRequest   = "ERR_INVALID_REQUEST"
	ErrNotFound         = "ERR_NOT_FOUND"
	ErrInternal         = "ERR_INTERNAL"
)

type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	ErrorCode  string      `json:"error_code,omitempty"`
}

func SendSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func SendError(c *fiber.Ctx, statusCode int, err error, errorCode string) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	} else {
		errMsg = errorCode
	}

	// Ghi log lỗi
	logger.WriteLogEx("error", errMsg, map[string]interface{}{
		"status_code": statusCode,
		"error_code":  errorCode,
		"path":        c.Path(),
		"method":      c.Method(),
	})

	return c.Status(statusCode).JSON(APIResponse{
		StatusCode: statusCode,
		Message:    errMsg,
		Error:      errMsg,
		ErrorCode:  errorCode,
	})
}

package middleware

import "net/http"

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Common error responses
var (
	ErrUnauthorized = ErrorResponse{
		Error:   "Unauthorized",
		Message: "Authentication required",
		Code:    "ERR_AUTH_REQUIRED",
	}

	ErrMissingAuthHeader = ErrorResponse{
		Error:   "Unauthorized",
		Message: "Missing authorization header",
		Code:    "ERR_AUTH_MISSING",
	}

	ErrInvalidAuthFormat = ErrorResponse{
		Error:   "Unauthorized",
		Message: "Invalid authorization header format",
		Code:    "ERR_AUTH_FORMAT",
	}

	ErrInvalidToken = ErrorResponse{
		Error:   "Unauthorized",
		Message: "Invalid or expired token",
		Code:    "ERR_AUTH_INVALID",
	}

	ErrForbidden = ErrorResponse{
		Error:   "Forbidden",
		Message: "Insufficient permissions",
		Code:    "ERR_PERMISSION_DENIED",
	}

	ErrRateLimitExceeded = ErrorResponse{
		Error:   "Rate limit exceeded",
		Message: "Too many requests. Please try again later.",
		Code:    "ERR_RATE_LIMIT",
	}
)

// StatusCode returns the appropriate HTTP status code for the error
func (e ErrorResponse) StatusCode() int {
	switch e.Code {
	case "ERR_AUTH_REQUIRED", "ERR_AUTH_MISSING", "ERR_AUTH_FORMAT", "ERR_AUTH_INVALID":
		return http.StatusUnauthorized
	case "ERR_PERMISSION_DENIED":
		return http.StatusForbidden
	case "ERR_RATE_LIMIT":
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

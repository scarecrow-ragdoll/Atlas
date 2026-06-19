package logger

import "github.com/google/uuid"

// HeaderRequestID is the canonical header name.
const HeaderRequestID = "X-Request-ID"

// generateID returns a new UUID v4 string.
func generateID() string {
	return uuid.NewString()
}

// extractID gets the request ID from the request header, or returns "".
func extractID(header string) string {
	return header
}

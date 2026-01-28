package utils

import (
	"os"
	"strings"

	config "car-license-number-fetcher/config"
)

// IsRequestFromMobile checks if the request is from a mobile device
func IsRequestFromMobile(userAgent string) bool {
	return strings.Contains(userAgent, config.MobileUserAgent)
}

// GetPort retrieves the port from environment variable or returns default
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = config.DefaultPort
	}
	return port
}

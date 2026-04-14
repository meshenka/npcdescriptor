package cmd

import "os"

// GetEnv returns environment variable value or default if empty.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

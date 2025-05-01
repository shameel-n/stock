package config

import (
	"os"
	"time"
)

// Config holds all configuration for the Upstox client
type Config struct {
	BaseURL    string
	APIKey     string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL:    "https://api.upstox.com/v2",
		APIKey:     os.Getenv("UPSTOX_API_KEY"),
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
	}
}

// NewConfig creates a new configuration with custom values
func NewConfig(baseURL, apiKey string, timeout time.Duration, maxRetries int, retryDelay time.Duration) *Config {
	return &Config{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Timeout:    timeout,
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}
}

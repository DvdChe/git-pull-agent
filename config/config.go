package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultSyncPath     = "/git-repo"
	defaultIntervalSec  = 60
	AuthMethodSSH       = "ssh"
	AuthMethodHTTPBasic = "http-basic"
	AuthMethodNone      = "none"
)

// Config holds the application's configuration parameters.
type Config struct {
	RepoURL          string
	SyncPath         string
	Interval         time.Duration
	AuthMethod       string
	SSHKeyPath       string
	SSHKeyPassphrase string
	HTTPUsername     string
	HTTPPassword     string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.RepoURL = os.Getenv("GIT_REPO_URL")
	if cfg.RepoURL == "" {
		return nil, fmt.Errorf("GIT_REPO_URL environment variable is required")
	}

	cfg.SyncPath = getEnvOrDefault("GIT_SYNC_PATH", defaultSyncPath)

	intervalStr := getEnvOrDefault("GIT_PULL_INTERVAL_SECONDS", strconv.Itoa(defaultIntervalSec))
	intervalSec, err := strconv.Atoi(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid value for GIT_PULL_INTERVAL_SECONDS: %w", err)
	}
	cfg.Interval = time.Duration(intervalSec) * time.Second

	cfg.AuthMethod = getEnvOrDefault("GIT_AUTH_METHOD", AuthMethodNone)

	switch cfg.AuthMethod {
	case AuthMethodSSH:
		cfg.SSHKeyPath = os.Getenv("GIT_SSH_KEY_PATH")
		if cfg.SSHKeyPath == "" {
			return nil, fmt.Errorf("GIT_SSH_KEY_PATH is required when GIT_AUTH_METHOD is ssh")
		}
		cfg.SSHKeyPassphrase = os.Getenv("GIT_SSH_KEY_PASSPHRASE")
	case AuthMethodHTTPBasic:
		cfg.HTTPUsername = os.Getenv("GIT_HTTP_USERNAME")
		if cfg.HTTPUsername == "" {
			return nil, fmt.Errorf("GIT_HTTP_USERNAME is required when GIT_AUTH_METHOD is http-basic")
		}
		cfg.HTTPPassword = os.Getenv("GIT_HTTP_PASSWORD")
		if cfg.HTTPPassword == "" {
			return nil, fmt.Errorf("GIT_HTTP_PASSWORD is required when GIT_AUTH_METHOD is http-basic")
		}
	case AuthMethodNone:
		// No specific auth config needed
	default:
		return nil, fmt.Errorf("unsupported GIT_AUTH_METHOD: %s", cfg.AuthMethod)
	}

	return cfg, nil
}

// getEnvOrDefault returns the value of the environment variable or a default value if not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

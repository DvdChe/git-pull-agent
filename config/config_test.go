package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

// Helper function to set environment variables for testing
func setEnv(key, value string) {
	os.Setenv(key, value)
}

// Helper function to unset environment variables after testing
func unsetEnv(key string) {
	os.Unsetenv(key)
}

// clearEnv clears all relevant environment variables
func clearEnv() {
	unsetEnv("GIT_REPO_URL")
	unsetEnv("GIT_SYNC_PATH")
	unsetEnv("GIT_PULL_INTERVAL_SECONDS")
	unsetEnv("GIT_AUTH_METHOD")
	unsetEnv("GIT_SSH_KEY_PATH")
	unsetEnv("GIT_SSH_KEY_PASSPHRASE")
	unsetEnv("GIT_HTTP_USERNAME")
	unsetEnv("GIT_HTTP_PASSWORD")
}

func TestLoadConfig_SuccessDefaults(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/test/repo.git")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.RepoURL != "https://github.com/test/repo.git" {
		t.Errorf("Expected RepoURL 'https://github.com/test/repo.git', got '%s'", cfg.RepoURL)
	}
	if cfg.SyncPath != defaultSyncPath {
		t.Errorf("Expected SyncPath '%s', got '%s'", defaultSyncPath, cfg.SyncPath)
	}
	if cfg.Interval != defaultIntervalSec*time.Second {
		t.Errorf("Expected Interval %s, got %s", defaultIntervalSec*time.Second, cfg.Interval)
	}
	if cfg.AuthMethod != AuthMethodNone {
		t.Errorf("Expected AuthMethod '%s', got '%s'", AuthMethodNone, cfg.AuthMethod)
	}
}

func TestLoadConfig_MissingRepoURL(t *testing.T) {
	clearEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing GIT_REPO_URL, got nil")
	}
	expectedError := "GIT_REPO_URL environment variable is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfig_CustomValues(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/custom/repo.git")
	setEnv("GIT_SYNC_PATH", "/tmp/custom_sync")
	setEnv("GIT_PULL_INTERVAL_SECONDS", "30")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.RepoURL != "https://github.com/custom/repo.git" {
		t.Errorf("Expected RepoURL 'https://github.com/custom/repo.git', got '%s'", cfg.RepoURL)
	}
	if cfg.SyncPath != "/tmp/custom_sync" {
		t.Errorf("Expected SyncPath '/tmp/custom_sync', got '%s'", cfg.SyncPath)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("Expected Interval 30s, got %s", cfg.Interval)
	}
}

func TestLoadConfig_SSHAuth(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "git@github.com:test/repo.git")
	setEnv("GIT_AUTH_METHOD", AuthMethodSSH)
	setEnv("GIT_SSH_KEY_PATH", "/path/to/ssh/key")
	setEnv("GIT_SSH_KEY_PASSPHRASE", "secret")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AuthMethod != AuthMethodSSH {
		t.Errorf("Expected AuthMethod '%s', got '%s'", AuthMethodSSH, cfg.AuthMethod)
	}
	if cfg.SSHKeyPath != "/path/to/ssh/key" {
		t.Errorf("Expected SSHKeyPath '/path/to/ssh/key', got '%s'", cfg.SSHKeyPath)
	}
	if cfg.SSHKeyPassphrase != "secret" {
		t.Errorf("Expected SSHKeyPassphrase 'secret', got '%s'", cfg.SSHKeyPassphrase)
	}
}

func TestLoadConfig_SSHAuthMissingKeyPath(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "git@github.com:test/repo.git")
	setEnv("GIT_AUTH_METHOD", AuthMethodSSH)
	// Missing GIT_SSH_KEY_PATH

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing GIT_SSH_KEY_PATH, got nil")
	}
	expectedError := "GIT_SSH_KEY_PATH is required when GIT_AUTH_METHOD is ssh"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfig_HTTPBasicAuth(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/test/repo.git")
	setEnv("GIT_AUTH_METHOD", AuthMethodHTTPBasic)
	setEnv("GIT_HTTP_USERNAME", "user")
	setEnv("GIT_HTTP_PASSWORD", "pass")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AuthMethod != AuthMethodHTTPBasic {
		t.Errorf("Expected AuthMethod '%s', got '%s'", AuthMethodHTTPBasic, cfg.AuthMethod)
	}
	if cfg.HTTPUsername != "user" {
		t.Errorf("Expected HTTPUsername 'user', got '%s'", cfg.HTTPUsername)
	}
	if cfg.HTTPPassword != "pass" {
		t.Errorf("Expected HTTPPassword 'pass', got '%s'", cfg.HTTPPassword)
	}
}

func TestLoadConfig_HTTPBasicAuthMissingUsername(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/test/repo.git")
	setEnv("GIT_AUTH_METHOD", AuthMethodHTTPBasic)
	// Missing GIT_HTTP_USERNAME
	setEnv("GIT_HTTP_PASSWORD", "pass")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing GIT_HTTP_USERNAME, got nil")
	}
	expectedError := "GIT_HTTP_USERNAME is required when GIT_AUTH_METHOD is http-basic"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfig_HTTPBasicAuthMissingPassword(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/test/repo.git")
	setEnv("GIT_AUTH_METHOD", AuthMethodHTTPBasic)
	setEnv("GIT_HTTP_USERNAME", "user")
	// Missing GIT_HTTP_PASSWORD

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing GIT_HTTP_PASSWORD, got nil")
	}
	expectedError := "GIT_HTTP_PASSWORD is required when GIT_AUTH_METHOD is http-basic"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfig_UnsupportedAuthMethod(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/test/repo.git")
	setEnv("GIT_AUTH_METHOD", "unsupported")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for unsupported GIT_AUTH_METHOD, got nil")
	}
	expectedError := "unsupported GIT_AUTH_METHOD: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfig_InvalidInterval(t *testing.T) {
	clearEnv()
	setEnv("GIT_REPO_URL", "https://github.com/test/repo.git")
	setEnv("GIT_PULL_INTERVAL_SECONDS", "not-a-number")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for invalid GIT_PULL_INTERVAL_SECONDS, got nil")
	}
	expectedErrorSubstring := "invalid value for GIT_PULL_INTERVAL_SECONDS"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorSubstring, err.Error())
	}
}
package test

import (
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"os"
	"testing"
)

func TestGetSecret(t *testing.T) {
	// Tests GetSecret function from config/secrets.go
	// GetSecret(key string) string
	// Tests if the function returns the value of the key from the environment variable
	t.Parallel()

	// Test if the function returns the value of the key from the environment variable
	t.Run("Get value from environment variable", func(t *testing.T) {
		key := "GITHUB_TOKEN"
		value := "token"
		// Set the environment variable
		_ = os.Setenv(key, value)

		// Call the function
		secret := config.GetSecret(key)

		// Check if the value is correct
		if secret != value {
			t.Errorf("Expected %s, got %s", value, secret)
		}

		// Unset the environment variable
		_ = os.Unsetenv(key)
	})
}

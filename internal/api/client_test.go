package api_test

import (
	"os"
	"testing"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

func unsetenv(t *testing.T, key string) {
	t.Helper()
	// t.Setenv registers the original value to be restored on cleanup.
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed to unset %s: %v", key, err)
	}
}

func TestNewClientFromPassedVariables(t *testing.T) {
	unsetenv(t, "COSMO_API_KEY")
	unsetenv(t, "COSMO_API_URL")

	client, err := api.NewClient("passed_api_key", "https://passed-url.com")
	if err != nil {
		t.Errorf("Expected client with passed variables, got error: %v", err)
	}

	if client.Client == nil {
		t.Errorf("Expected client to be created but got nil")
	}
}

func TestNewClientFromEnvironment(t *testing.T) {
	t.Setenv("COSMO_API_KEY", "env_api_key")
	t.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err != nil {
		t.Errorf("Expected client with env variables, got error: %v", err)
	}

	if client.Client == nil {
		t.Errorf("Expected client to be created but got nil")
	}
}

func TestNewClientFromEnvironmentWithoutApiKey(t *testing.T) {
	unsetenv(t, "COSMO_API_KEY")
	t.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err == nil {
		t.Errorf("Expected client creation to fail but got client: %v", err)
	}

	if client != nil {
		t.Errorf("Expected client not to be created")
	}
}

func TestNewClientFromEnvironmentWithoutApiUrlAndApiKey(t *testing.T) {
	unsetenv(t, "COSMO_API_KEY")
	unsetenv(t, "COSMO_API_URL")

	client, err := api.NewClient("", "")
	if err == nil {
		t.Errorf("Expected client not to be created: %v", err)
	}

	if client != nil {
		t.Errorf("Expected client not to be created")
	}
}

func TestNewClientFromEnvironmentWithApiKey(t *testing.T) {
	unsetenv(t, "COSMO_API_KEY")
	unsetenv(t, "COSMO_API_URL")

	client, err := api.NewClient("cosmo_api_key", "")
	if err != nil {
		t.Errorf("Expected client to be created but got error: %v", err)
	}

	if client.Client == nil {
		t.Errorf("Expected client to be created but got nil")
	}
}

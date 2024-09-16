package api_test

import (
	"os"
	"testing"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

func TestNewClientFromPassedVariables(t *testing.T) {
	os.Setenv("COSMO_API_KEY", "")
	os.Setenv("COSMO_API_URL", "")

	client, err := api.NewClient("passed_api_key", "https://passed-url.com")
	if err != nil || client.Client == nil {
		t.Errorf("Expected client with passed variables, got error: %v", err)
	}

	os.Unsetenv("COSMO_API_KEY")
	os.Unsetenv("COSMO_API_URL")
}

func TestNewClientFromEnvironment(t *testing.T) {
	os.Setenv("COSMO_API_KEY", "env_api_key")
	os.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err != nil || client.Client == nil {
		t.Errorf("Expected client with env variables, got error: %v", err)
	}

	os.Unsetenv("COSMO_API_KEY")
	os.Unsetenv("COSMO_API_URL")
}

func TestNewClientFromEnvironmentWithoutApiKey(t *testing.T) {
	os.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err != nil && client != nil {
		t.Errorf("Expected client creation to fail but got client: %v", err)
	}

	os.Unsetenv("COSMO_API_URL")
}

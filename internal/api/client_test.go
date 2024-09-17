package api_test

import (
	"os"
	"testing"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

func TestNewClientFromPassedVariables(t *testing.T) {
	os.Unsetenv("COSMO_API_KEY")
	os.Unsetenv("COSMO_API_URL")

	client, err := api.NewClient("passed_api_key", "https://passed-url.com")
	if err != nil {
		t.Errorf("Expected client with passed variables, got error: %v", err)
	}

	if client.Client == nil {
		t.Errorf("Expected client to be created but got nil")
	}
}

func TestNewClientFromEnvironment(t *testing.T) {
	os.Setenv("COSMO_API_KEY", "env_api_key")
	os.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err != nil {
		t.Errorf("Expected client with env variables, got error: %v", err)
	}

	if client.Client == nil {
		t.Errorf("Expected client to be created but got nil")
	}

	os.Unsetenv("COSMO_API_KEY")
	os.Unsetenv("COSMO_API_URL")
}

func TestNewClientFromEnvironmentWithoutApiKey(t *testing.T) {
	os.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err == nil {
		t.Errorf("Expected client creation to fail but got client: %v", err)
	}

	if client != nil {
		t.Errorf("Expected client not to be created")
	}

	os.Unsetenv("COSMO_API_URL")
}

func TestNewClientFromEnvironmentWithoutApiUrlAndApiKey(t *testing.T) {
	os.Unsetenv("COSMO_API_KEY")
	os.Unsetenv("COSMO_API_URL")

	client, err := api.NewClient("", "")
	if err == nil {
		t.Errorf("Expected client not to be created: %v", err)
	}

	if client != nil {
		t.Errorf("Expected client not to be created")
	}
}

func TestNewClientFromEnvironmentWithApiKey(t *testing.T) {
	os.Unsetenv("COSMO_API_KEY")
	os.Unsetenv("COSMO_API_URL")

	client, err := api.NewClient("cosmo_api_key", "")
	if err != nil {
		t.Errorf("Expected client to be created but got error: %v", err)
	}

	if client.Client == nil {
		t.Errorf("Expected client to be created but got nil")
	}
}

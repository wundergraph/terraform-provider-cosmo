package api_test

import (
	"testing"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

func TestNewClient(t *testing.T) {
	t.Setenv("COSMO_API_KEY", "env_api_key")
	t.Setenv("COSMO_API_URL", "https://env-url.com")

	client, err := api.NewClient("", "")
	if err != nil || client.CosmoApiKey != "env_api_key" || client.Client == nil {
		t.Errorf("Expected client with env variables, got error: %v", err)
	}

	t.Setenv("COSMO_API_KEY", "")
	t.Setenv("COSMO_API_URL", "")

	client, err = api.NewClient("passed_api_key", "https://passed-url.com")
	if err != nil || client.CosmoApiKey != "passed_api_key" || client.Client == nil {
		t.Errorf("Expected client with passed variables, got error: %v", err)
	}

	client, err = api.NewClient("", "")
	if err == nil {
		t.Errorf("Expected error when no API key is provided, got client: %v", client)
	}
}

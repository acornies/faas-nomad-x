package handlers

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/hashicorp/vault/api"
	vhttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

type SecretsTestConfig struct {
	Method          string
	URL             string
	Body            []byte
	Config          types.VaultConfig
	SecretsResponse SecretsResponse
	SecretsError    error
}

// Unapologetically taken from:
// https://stackoverflow.com/questions/57771228/mocking-hashicorp-vault-in-go
func createTestVault(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := vhttp.TestServer(t, core)
	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)

	// Setup required secrets, policies, etc.
	_, err = client.Logical().Write("secret/openfaas", map[string]interface{}{
		"secret": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}

func setupSecrets(t *testing.T, secrets SecretsTestConfig) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {

	return MakeSecrets(&secrets.Config), httptest.NewRecorder(),
		httptest.NewRequest(secrets.Method, secrets.URL, bytes.NewReader(secrets.Body))
}

func TestListSecrets(t *testing.T) {

	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	// memory core Vault default secret backend
	providerConfig.Vault.Secrets.KeyPrefix = "secret/openfaas"
	providerConfig.Vault.Secrets.KVVersion = 1
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "GET",
		URL:    "/system/secrets"})

	h(recorder, r)

	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

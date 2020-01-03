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
	Method string
	URL    string
	Body   []byte
	Config types.VaultConfig
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
	_, err = client.Logical().Write("secret/openfaas/cows", map[string]interface{}{
		"value": "changeme",
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}

func setupSecrets(t *testing.T, secrets SecretsTestConfig) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {

	// memory core Vault default secret backend
	secrets.Config.Secrets.KeyPrefix = "secret/openfaas"
	secrets.Config.Secrets.KVVersion = 1

	return MakeSecrets(&secrets.Config), httptest.NewRecorder(),
		httptest.NewRequest(secrets.Method, secrets.URL, bytes.NewReader(secrets.Body))
}

func TestListSecrets(t *testing.T) {

	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()

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

func TestCreateSecret(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "POST",
		URL:    "/system/secrets",
		Body: []byte(`
		{
			"name": "figlet",
			"value": "changeme"
		}`)})

	h(recorder, r)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestUpdateSecret(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "PUT",
		URL:    "/system/secrets",
		Body: []byte(`
		{
			"name": "cows",
			"value": "changed"
		}`)})

	h(recorder, r)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestDeleteSecret(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "DELETE",
		URL:    "/system/secrets",
		Body: []byte(`
		{
			"name": "cows"
		}`)})

	h(recorder, r)

	if recorder.Code != http.StatusNoContent {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestBadRequest(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "BLAH",
		URL:    "/system/secrets",
		Body:   []byte(``)})

	h(recorder, r)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestFailedVaultRequest(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	ln.Close() // intentionally shut down

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "GET",
		URL:    "/system/secrets"})

	h(recorder, r)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestBadCreateUpdateSecret(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "POST",
		URL:    "/system/secrets",
		Body:   []byte(``)})

	h(recorder, r)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestBadDeleteSecret(t *testing.T) {
	ln, client := createTestVault(t)
	providerConfig := types.NewProviderConfig().Default()
	providerConfig.Vault.Client = client
	defer ln.Close()

	h, recorder, r := setupSecrets(t, SecretsTestConfig{
		Config: providerConfig.Vault,
		Method: "DELETE",
		URL:    "/system/secrets",
		Body:   []byte(``)})

	h(recorder, r)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

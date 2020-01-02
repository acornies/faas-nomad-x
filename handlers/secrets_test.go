package handlers

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/vault"
)

type SecretsTestConfig struct {
	Method          string
	URL             string
	Body            []byte
	SecretsResponse SecretsResponse
	SecretsError    error
}

// Unapologetically taken from:
// https://stackoverflow.com/questions/57771228/mocking-hashicorp-vault-in-go
func createTestVault(t *testing.T, prefix string) (net.Listener, *api.Client) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

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
	_, err = client.Logical().Write(prefix, map[string]interface{}{
		"secret": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}

func setupSecrets(secrets SecretsTestConfig) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {

	testConfig := types.NewVaultConfig()

	ln, client := createTestVault(t, testConfig.Secrets.KeyPrefix)
	defer ln.Close()

	testConfig.Client = &client

	return MakeSecrets(testConfig), httptest.NewRecorder(),
		httptest.NewRequest(secrets.Method, secrets.URL, bytes.NewReader(secrets.Body))
}

func TestListSecrets(t *testing.T) {

	h, recorder, r := setupSecrets(SecretsTestConfig{
		Method: "GET",
		URL:    "/system/secrets"})

	h(recorder, r)
}

package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	xtest "github.com/acornies/faas-nomad-x/testing"
	"github.com/acornies/faas-nomad-x/types"
)

type SecretsTestConfig struct {
	Method string
	URL    string
	Body   []byte
	// RegisterResponse api.JobRegisterResponse
	// RegisterError    error
}

func setupSecrets(secrets SecretsTestConfig) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {
	testConfig := types.ProviderConfig{}
	testConfig.Default()

	mockClient := xtest.MockVaultClient{}

	return MakeSecrets(&testConfig, &mockClient), httptest.NewRecorder(),
		httptest.NewRequest(secrets.Method, secrets.URL, bytes.NewReader(secrets.Body))
}

package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acornies/faas-nomad-x/types"
)

func setup(method string, url string, body []byte) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {
	testConfig := types.ProviderConfig{}
	testConfig.Default()
	return MakeDeploy(&testConfig), httptest.NewRecorder(),
		httptest.NewRequest(method, url, bytes.NewReader(body))
}

func TestDeployEmptyBody(t *testing.T) {

	h, recorder, r := setup("POST", "/system/functions", []byte(""))

	h(recorder, r)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

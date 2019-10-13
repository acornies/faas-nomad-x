package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/mock"
)

func setup(method string, url string, body []byte) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {
	testConfig := types.ProviderConfig{}
	testConfig.Default()
	mockJobs := types.MockJobs{}

	mockJobs.On("Register", mock.Anything, mock.Anything).Return(
		&api.JobRegisterResponse{JobModifyIndex: 1},
		nil, nil)

	return MakeDeploy(&testConfig, &mockJobs), httptest.NewRecorder(),
		httptest.NewRequest(method, url, bytes.NewReader(body))
}

func TestDeployEmptyBody(t *testing.T) {

	h, recorder, r := setup("POST", "/system/functions", []byte(""))

	h(recorder, r)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestDeployAllDefaultConfig(t *testing.T) {

	h, recorder, r := setup("POST", "/system/functions", []byte(`
	{
		"service": "nodeinfo",
		"network": "func_functions",
		"image": "functions/nodeinfo:latest",
		"envProcess": "node main.js",
		"envVars": {
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string"
		},
		"constraints": [
			"node.platform.os == linux"
		],
		"labels": {
			"foo": "bar"
		},
		"annotations": {
			"topics": "awesome-kafka-topic",
			"foo": "bar"
		},
		"secrets": [
			"secret-name-1"
		],
		"registryAuth": "dXNlcjpwYXNzd29yZA==",
		"limits": {
			"memory": "128M",
			"cpu": "0.01"
		},
		"requests": {
			"memory": "128M",
			"cpu": "0.01"
		},
		"readOnlyRootFilesystem": true
	}`))

	h(recorder, r)

	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

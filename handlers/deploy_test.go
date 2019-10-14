package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/mock"
)

type DeployTestConfig struct {
	Method           string
	URL              string
	Body             []byte
	RegisterResponse api.JobRegisterResponse
	RegisterError    error
}

func setup(deploy DeployTestConfig) (http.HandlerFunc, *httptest.ResponseRecorder, *http.Request) {
	testConfig := types.ProviderConfig{}
	testConfig.Default()

	mockJobs := types.MockJobs{}

	mockJobs.On("Register", mock.Anything, mock.Anything).Return(
		&deploy.RegisterResponse,
		nil, deploy.RegisterError)

	return MakeDeploy(&testConfig, &mockJobs), httptest.NewRecorder(),
		httptest.NewRequest(deploy.Method, deploy.URL, bytes.NewReader(deploy.Body))
}

func TestDeployEmptyBody(t *testing.T) {

	h, recorder, r := setup(DeployTestConfig{
		Method: "POST",
		URL:    "/system/functions",
		Body:   []byte("")})

	h(recorder, r)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestDeployRegisterError(t *testing.T) {

	h, recorder, r := setup(DeployTestConfig{
		Method: "POST",
		URL:    "/system/functions",
		Body: []byte(`
		{
			"service": "nodeinfo",
			"network": "func_functions",
			"image": "functions/nodeinfo:latest",
			"envProcess": "node main.js"
		}`),
		RegisterResponse: api.JobRegisterResponse{JobModifyIndex: 1},
		RegisterError:    errors.New("Register() returns error"),
	})

	h(recorder, r)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

func TestDeployAllDefaultConfig(t *testing.T) {

	h, recorder, r := setup(DeployTestConfig{
		Method:           "POST",
		URL:              "/system/functions",
		RegisterResponse: api.JobRegisterResponse{JobModifyIndex: 1},
		RegisterError:    nil,
		Body: []byte(`
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
	}`)})

	h(recorder, r)

	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected response code %d", recorder.Code)
	}
}

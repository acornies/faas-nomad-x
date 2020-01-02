package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/hashicorp/vault/api"
	providerTypes "github.com/openfaas/faas-provider/types"
)

var (
	response SecretsResponse
)

type SecretsResponse struct {
	StatusCode int
	Body       []byte
}

func MakeSecrets(config *types.VaultConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println("Couldn't read body of a request", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodGet:
			response, err = getSecrets(config, body)
			break
			// case http.MethodPost:
			// 	response, err = createNewSecret(vs, http.MethodPost, body)
			// 	break
			// case http.MethodPut:
			// 	response, err = createNewSecret(vs, http.MethodPut, body)
			// 	break
			// case http.MethodDelete:
			// 	response, err = deleteSecret(vs, body)
			// 	break
		}

		if err != nil {
			log.Println("Error in secrets handler", err)
			w.WriteHeader(response.StatusCode)
			return
		}

		w.WriteHeader(response.StatusCode)

		if response.Body != nil {
			_, err := w.Write(response.Body)

			if err != nil {
				log.Println("Cannot write body of a response")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func getSecrets(vc *types.VaultConfig, body []byte) (resp SecretsResponse, err error) {

	response, err := vc.DoRequest("LIST",
		fmt.Sprintf("/v1/%s", vc.Secrets.KeyPrefix), nil)

	if err != nil {
		return SecretsResponse{StatusCode: http.StatusInternalServerError},
			fmt.Errorf("Error reading response body: %s", err)
	}

	// If Vault finds nothing, return StatusOK according to gateway API docs
	if response.StatusCode == http.StatusNotFound {
		return SecretsResponse{
			StatusCode: http.StatusOK,
			Body:       []byte(`[]`),
		}, nil
	}

	var secretList api.Secret
	secretsBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return SecretsResponse{StatusCode: http.StatusInternalServerError},
			fmt.Errorf("Error reading response body: %s", err)
	}

	err = json.Unmarshal(secretsBody, &secretList)
	if err != nil {
		return SecretsResponse{StatusCode: http.StatusInternalServerError},
			fmt.Errorf("Error in json deserialisation: %s", err)
	}

	secrets := []providerTypes.Secret{}
	for _, k := range secretList.Data["keys"].([]interface{}) {
		secrets = append(secrets, providerTypes.Secret{Name: k.(string)})
	}

	resultsJson, _ := json.Marshal(secrets)

	return SecretsResponse{StatusCode: http.StatusOK, Body: resultsJson}, nil
}

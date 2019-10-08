package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/openfaas/faas/gateway/requests"
)

var ()

func MakeDeploy(config *types.ProviderConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		req := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

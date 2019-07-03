package handlers

import (
	"net/http"

	"github.com/acornies/faas-nomad-x/types"
)

func MakeDeploy(config *types.ProviderConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		return
	}
}

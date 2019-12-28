package handlers

import (
	"net/http"

	"github.com/acornies/faas-nomad-x/types"
)

func MakeSecrets(config *types.ProviderConfig, client types.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}
}

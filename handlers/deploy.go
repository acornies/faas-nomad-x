package handlers

import (
	"net/http"

	"github.com/acornies/faas-nomad-x/types"
	nomadapi "github.com/hashicorp/nomad/api"
)

func MakeDeploy(config *types.ProviderConfig, nomadClient *nomadapi.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		return
	}
}

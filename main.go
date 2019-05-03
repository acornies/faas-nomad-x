package main

import (
	"flag"
	"log"
	"net/http"

	provider "github.com/acornies/faas-nomad-x/types"
	bootstrap "github.com/openfaas/faas-provider"
	"github.com/openfaas/faas-provider/types"
)

var (
	// listenPort = flag.Int("listen-port", 8081, "Port to bind the server to")
	configFile = flag.String("config-file", "./config.d/config.hcl", "The provider configuration directoy. One or many HCL or JSON files.")
)

func main() {

	file := *configFile
	providerConfig, err := provider.NewProviderConfig().LoadFile(file)
	if err != nil {
		log.Print("Failed to load config-dir", err)
	}

	faasConfig := &types.FaaSConfig{
		TCPPort: &providerConfig.ListenPort,
	}

	handlers := &types.FaaSHandlers{

		FunctionReader: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		DeployHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		FunctionProxy: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		ReplicaReader: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		ReplicaUpdater: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		SecretHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		DeleteHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		Health: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		InfoHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},
	}

	bootstrap.Serve(handlers, faasConfig)
}

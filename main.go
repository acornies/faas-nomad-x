package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/acornies/faas-nomad-x/types"
	bootstrap "github.com/openfaas/faas-provider"
	btypes "github.com/openfaas/faas-provider/types"
)

var (
	listenPort = flag.Int("listen-port", 0, "Server listen port override")
	consulAddr = flag.String("consul-addr", "", "Consul address override")
	nomadAddr  = flag.String("nomad-addr", "", "Nomad address override")
	vaultAddr  = flag.String("vault-addr", "", "Vault address override")
	configFile = flag.String("config-file", "./default.hcl", "The provider configuration file. Either HCL or JSON format.")
)

func main() {

	file := *configFile
	providerConfig, err := types.NewProviderConfig().LoadFile(file)
	if err != nil {
		log.Fatal("Failed to load config-file", err)
	}

	flag.Parse()
	port := *listenPort
	consul := *consulAddr
	nomad := *nomadAddr
	vault := *vaultAddr

	providerConfig.LoadCommandLine(port, consul, nomad, vault)

	faasConfig := &btypes.FaaSConfig{
		TCPPort:         &providerConfig.ListenPort,
		EnableBasicAuth: providerConfig.AuthEnabled,
		SecretMountPath: providerConfig.CredentialsDir,
	}

	handlers := &btypes.FaaSHandlers{

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

		HealthHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		InfoHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},
	}

	bootstrap.Serve(handlers, faasConfig)
}

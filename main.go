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
	listenPort = flag.Int("listen-port", 8080, "Server listen port override")
	consulAddr = flag.String("consul-addr", "", "Consul address override")
	nomadAddr  = flag.String("nomad-addr", "", "Nomad address override")
	vaultAddr  = flag.String("vault-addr", "", "Vault address override")
	configFile = flag.String("config-file", "", "The provider configuration file. Either HCL or JSON format.")
)

func main() {

	flag.Parse()
	file := *configFile

	port := *listenPort
	consul := *consulAddr
	nomad := *nomadAddr
	vault := *vaultAddr

	providerConfig, err := configure(file)
	if err != nil {
		log.Printf("Error loading config file: %v. Using defaults...", err)
	}
	providerConfig.LoadCommandLine(port, consul, nomad, vault)

	faasConfig := &btypes.FaaSConfig{
		TCPPort:         &providerConfig.ListenPort,
		EnableBasicAuth: providerConfig.Auth.Enabled,
		SecretMountPath: providerConfig.Auth.CredentialsDir,
		ReadTimeout:     providerConfig.ReadTimeout,
		WriteTimeout:    providerConfig.WriteTimeout,
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

		LogHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},

		InfoHandler: func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
		},
	}

	bootstrap.Serve(handlers, faasConfig)
}

func configure(file string) (*types.ProviderConfig, error) {
	config := types.NewProviderConfig()
	if len(file) == 0 {
		log.Print("No configuration file detected. Using defaults...")
		config.Default()
		return config, nil
	} else {
		config, err := config.LoadFile(file)
		return config, err
	}
}

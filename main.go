package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/acornies/faas-nomad-x/handlers"
	"github.com/acornies/faas-nomad-x/types"
	bootstrap "github.com/openfaas/faas-provider"
	btypes "github.com/openfaas/faas-provider/types"
)

var (
	listenPort    = flag.Int(types.ListenPort, 8080, "Server listen port override")
	consulAddr    = flag.String(types.ConsulAddr, "", "Consul address override")
	consulTLSSkip = flag.Bool(types.ConsulTLSSkipVerify, false, "Skip TLS verification for Consul address")
	nomadAddr     = flag.String(types.NomadAddr, "", "Nomad address override")
	nomadTLSSkip  = flag.Bool(types.NomadTLSSkipVerify, false, "Skip TLS verification for Nomad address")
	vaultAddr     = flag.String(types.VaultAddr, "", "Vault address override")
	vaultTLSSkip  = flag.Bool(types.VaultTLSSkipVerify, false, "Skip TLS verification for Vault address")
	configFile    = flag.String("config-file", "", "The provider configuration file. Either HCL or JSON format.")
)

func main() {

	flag.Parse()

	port := *listenPort
	consul := *consulAddr
	consulSkip := *consulTLSSkip
	nomad := *nomadAddr
	nomadSkip := *nomadTLSSkip
	vault := *vaultAddr
	vaultSkip := *vaultTLSSkip

	override := map[string]interface{}{
		types.ListenPort:          port,
		types.ConsulAddr:          consul,
		types.ConsulTLSSkipVerify: consulSkip,
		types.NomadAddr:           nomad,
		types.NomadTLSSkipVerify:  nomadSkip,
		types.VaultAddr:           vault,
		types.VaultTLSSkipVerify:  vaultSkip,
	}

	providerConfig, err := configure(configFile)
	if err != nil {
		log.Printf("Error loading config file: %v. Using defaults...", err)
	}
	providerConfig.LoadCommandLine(override)

	var nomadClient types.Jobs
	nomadClient, err = providerConfig.Nomad.MakeClient()
	if err != nil {
		log.Fatal("Failed to create Nomad client ", err)
	}

	var vaultClient types.Client
	vaultClient, err = providerConfig.Vault.MakeClient()
	if err != nil {
		log.Print("WARN: Failed to create Vault client ", err)
	}

	faasConfig := &btypes.FaaSConfig{
		TCPPort:         &providerConfig.ListenPort,
		EnableBasicAuth: providerConfig.Auth.Enabled,
		SecretMountPath: providerConfig.Auth.CredentialsDir,
		ReadTimeout:     providerConfig.ReadTimeout,
		WriteTimeout:    providerConfig.WriteTimeout,
	}

	handlers := &btypes.FaaSHandlers{

		FunctionReader: func(w http.ResponseWriter, r *http.Request) {},

		DeployHandler: handlers.MakeDeploy(providerConfig, nomadClient),

		FunctionProxy: func(w http.ResponseWriter, r *http.Request) {},

		ReplicaReader: func(w http.ResponseWriter, r *http.Request) {},

		ReplicaUpdater: func(w http.ResponseWriter, r *http.Request) {},

		SecretHandler: handlers.MakeSecrets(providerConfig, vaultClient),

		DeleteHandler: func(w http.ResponseWriter, r *http.Request) {},

		HealthHandler: func(w http.ResponseWriter, r *http.Request) {},

		LogHandler: func(w http.ResponseWriter, r *http.Request) {},

		InfoHandler: func(w http.ResponseWriter, r *http.Request) {},
	}

	bootstrap.Serve(handlers, faasConfig)
}

func configure(configFile *string) (*types.ProviderConfig, error) {
	config := types.NewProviderConfig()
	file := *configFile
	if len(file) == 0 {
		log.Print("No configuration file detected. Using defaults...")
		config.Default()
		return config, nil
	} else {
		config, err := config.LoadFile(file)
		return config, err
	}
}

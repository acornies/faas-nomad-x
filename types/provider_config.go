package types

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type ProviderConfig struct {
	LogLevel       string `hcl:"log_level"`
	ListenPort     int    `hcl:"listen_port"`
	HealthEnabled  bool   `hcl:"health_enabled"`
	AuthEnabled    bool   `hcl:"auth_enabled"`
	CredentialsDir string `hcl:"credentials_dir"`
	Nomad          NomadConfig
	Consul         ConsulConfig
	Vault          VaultConfig
}

type NomadConfig struct {
	Address  string
	ACLToken string `hcl:"acl_token"`
	TLS      TLSConfig
	Region   string
	Driver   string
}

type VaultConfig struct {
	Address string
	TLS     TLSConfig
	AppRole AppRoleConfig
	Secrets SecretConfig
}

type AppRoleConfig struct {
	RoleID   string `hcl:"role_id"`
	SecretID string `hcl:"secret_id"`
}

type SecretConfig struct {
	KeyPrefix string `hcl:"key_prefix"`
	KVVersion int    `hcl:"kv_version"`
	Policy    string
}

type ConsulConfig struct {
	Address    string
	ACLToken   string `hcl:"acl_token"`
	TLS        TLSConfig
	DNSEnabled bool `hcl:"dns_enabled"`
}

type TLSConfig struct {
	Insecure bool
	CAFile   string `hcl:"ca_file"`
	CertFile string `hcl:"cert_file"`
	KeyFile  string `hcl:"key_file"`
}

func NewProviderConfig() *ProviderConfig {
	pc := &ProviderConfig{}
	return pc
}

func (pc *ProviderConfig) Default() *ProviderConfig {
	pc.ListenPort = 8081
	pc.Consul = ConsulConfig{
		Address: "127.0.0.1:8500",
	}
	pc.Nomad = NomadConfig{
		Region:  "global",
		Address: "127.0.0.1:4646",
	}
	pc.Vault = VaultConfig{
		Address: "127.0.0.1:8200",
		Secrets: SecretConfig{
			KeyPrefix: "kv/openfaas",
			KVVersion: 2,
			Policy:    "openfaas",
		},
	}
	return pc
}

func (pc *ProviderConfig) LoadFile(configFile string) (*ProviderConfig, error) {
	pc = pc.Default()
	fileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	decoded := &ProviderConfig{}
	err = hcl.Unmarshal(fileBytes, decoded)
	if err != nil {
		return nil, err
	}

	pc.ListenPort = intOrDefault(decoded.ListenPort, pc.ListenPort)
	pc.Consul.ACLToken = stringOrDefault(decoded.Consul.ACLToken, pc.Consul.ACLToken)
	pc.Consul.Address = stringOrDefault(decoded.Consul.Address, pc.Consul.Address)
	pc.Consul.DNSEnabled = decoded.Consul.DNSEnabled
	pc.Consul.TLS = decoded.Consul.TLS
	// TODO: continue

	return pc, nil
}

func (pc *ProviderConfig) LoadCommandLine(listenPort int, consulAddr, nomadAddr, vaultAddr string) *ProviderConfig {
	pc.ListenPort = listenPort
	pc.Consul.Address = consulAddr
	pc.Nomad.Address = nomadAddr
	pc.Vault.Address = vaultAddr
	return pc
}

func stringOrDefault(value, fallback string) string {
	if len(value) <= 0 {
		return fallback
	}
	return value
}

func intOrDefault(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

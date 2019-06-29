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
	DNSServers     bool   `hcl:"dns_servers"`
	Nomad          NomadConfig
	Consul         ConsulConfig
	Vault          VaultConfig
}

type NomadConfig struct {
	Address  string    `hcl:"address"`
	ACLToken string    `hcl:"acl_token"`
	TLS      TLSConfig `hcl:"tls"`
	Region   string    `hcl:"region"`
	Driver   string    `hcl:"driver"`
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
	Address  string
	ACLToken string `hcl:"acl_token"`
	TLS      TLSConfig
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

	pc.LogLevel = stringOrDefault(decoded.LogLevel, pc.LogLevel)
	pc.ListenPort = intOrDefault(decoded.ListenPort, pc.ListenPort)
	pc.HealthEnabled = decoded.HealthEnabled
	pc.AuthEnabled = decoded.AuthEnabled
	pc.CredentialsDir = stringOrDefault(decoded.CredentialsDir, pc.CredentialsDir)

	pc.Consul.ACLToken = stringOrDefault(decoded.Consul.ACLToken, pc.Consul.ACLToken)
	pc.Consul.Address = stringOrDefault(decoded.Consul.Address, pc.Consul.Address)
	pc.Consul.TLS = decoded.Consul.TLS

	pc.Nomad.Region = stringOrDefault(decoded.Nomad.Region, pc.Nomad.Region)
	pc.Nomad.Driver = stringOrDefault(decoded.Nomad.Driver, pc.Nomad.Driver)
	pc.Nomad.ACLToken = stringOrDefault(decoded.Nomad.ACLToken, pc.Nomad.ACLToken)
	pc.Nomad.Address = stringOrDefault(decoded.Nomad.Address, pc.Nomad.Address)
	pc.Nomad.TLS = decoded.Nomad.TLS

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

package types

import (
	"io/ioutil"
	"time"

	"github.com/hashicorp/hcl"
	nomadapi "github.com/hashicorp/nomad/api"
)

type ProviderConfig struct {
	LogLevel      string `hcl:"log_level"`
	ListenPort    int    `hcl:"listen_port"`
	HealthEnabled bool   `hcl:"health_enabled"`
	Auth          AuthConfig
	DNSServers    bool `hcl:"dns_servers"`
	Nomad         NomadConfig
	Consul        ConsulConfig
	Vault         VaultConfig
	ReadTimeout   time.Duration `hcl:"read_timeout"`
	WriteTimeout  time.Duration `hcl:"write_timeout"`
}

type AuthConfig struct {
	Enabled        bool
	Type           string
	CredentialsDir string `hcl:"credentials_dir"`
}

type NomadConfig struct {
	Client   *nomadapi.Client
	Address  string `hcl:"address"`
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
	KeyPrefix  string `hcl:"key_prefix"`
	KVVersion  int    `hcl:"kv_version"`
	PolicyName string `hcl:"policy_name"`
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
	pc.LogLevel = "INFO"
	pc.Consul = ConsulConfig{
		Address: "127.0.0.1:8500",
	}
	pc.Nomad = NomadConfig{
		Region:  "global",
		Address: "http://127.0.0.1:4646",
		Driver:  "docker",
	}
	pc.Vault = VaultConfig{
		Address: "127.0.0.1:8200",
		Secrets: SecretConfig{
			KeyPrefix:  "kv/openfaas",
			KVVersion:  2,
			PolicyName: "openfaas",
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
	pc.Auth.Enabled = decoded.Auth.Enabled
	pc.Auth.CredentialsDir = stringOrDefault(decoded.Auth.CredentialsDir, pc.Auth.CredentialsDir)

	pc.Consul.ACLToken = stringOrDefault(decoded.Consul.ACLToken, pc.Consul.ACLToken)
	pc.Consul.Address = stringOrDefault(decoded.Consul.Address, pc.Consul.Address)
	pc.Consul.TLS = decoded.Consul.TLS

	pc.Nomad.Region = stringOrDefault(decoded.Nomad.Region, pc.Nomad.Region)
	pc.Nomad.Driver = stringOrDefault(decoded.Nomad.Driver, pc.Nomad.Driver)
	pc.Nomad.ACLToken = stringOrDefault(decoded.Nomad.ACLToken, pc.Nomad.ACLToken)
	pc.Nomad.Address = stringOrDefault(decoded.Nomad.Address, pc.Nomad.Address)
	pc.Nomad.TLS = decoded.Nomad.TLS

	pc.Vault.Address = stringOrDefault(decoded.Vault.Address, pc.Vault.Address)
	pc.Vault.AppRole.RoleID = stringOrDefault(decoded.Vault.AppRole.RoleID, pc.Vault.AppRole.RoleID)
	pc.Vault.AppRole.SecretID = stringOrDefault(decoded.Vault.AppRole.RoleID, pc.Vault.AppRole.SecretID)
	pc.Vault.Secrets.KVVersion = intOrDefault(decoded.Vault.Secrets.KVVersion, pc.Vault.Secrets.KVVersion)
	pc.Vault.Secrets.KeyPrefix = stringOrDefault(decoded.Vault.Secrets.KeyPrefix, pc.Vault.Secrets.KeyPrefix)
	pc.Vault.Secrets.PolicyName = stringOrDefault(decoded.Vault.Secrets.PolicyName, pc.Vault.Secrets.PolicyName)
	pc.Vault.TLS = decoded.Vault.TLS
	return pc, nil
}

func (pc *ProviderConfig) LoadCommandLine(listenPort int, consulAddr, nomadAddr, vaultAddr string) *ProviderConfig {
	pc.ListenPort = intOrDefault(listenPort, pc.ListenPort)
	pc.Consul.Address = stringOrDefault(consulAddr, pc.Consul.Address)
	pc.Nomad.Address = stringOrDefault(nomadAddr, pc.Nomad.Address)
	pc.Vault.Address = stringOrDefault(vaultAddr, pc.Nomad.Address)
	return pc
}

func (pc *ProviderConfig) MakeNomadClient() error {
	client, err := nomadapi.NewClient(&nomadapi.Config{
		Address:  pc.Nomad.Address,
		Region:   pc.Nomad.Region,
		SecretID: pc.Nomad.ACLToken,
		TLSConfig: &nomadapi.TLSConfig{
			CACert:     pc.Nomad.TLS.CAFile,
			ClientCert: pc.Nomad.TLS.CertFile,
			ClientKey:  pc.Nomad.TLS.KeyFile,
			Insecure:   pc.Nomad.TLS.Insecure,
		},
	})
	if err != nil {
		return err
	} else {
		pc.Nomad.Client = client
		return nil
	}
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

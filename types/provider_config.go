package types

import (
	"io/ioutil"
	"time"

	"github.com/hashicorp/hcl"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/imdario/mergo"
)

const ListenPort = "listen-port"
const ConsulAddr = "consul-addr"
const ConsulTLSSkipVerify = "consul-tls-skip-verify"
const NomadAddr = "nomad-addr"
const NomadTLSSkipVerify = "nomad-tls-skip-verify"
const VaultAddr = "vault-addr"
const VaultTLSSkipVerify = "vault-tls-skip-verify"

type ProviderConfig struct {
	LogLevel       string `hcl:"log_level"`
	ListenPort     int    `hcl:"listen_port"`
	HealthEnabled  bool   `hcl:"health_enabled"`
	Auth           AuthConfig
	DNSServers     bool `hcl:"dns_servers"`
	Nomad          NomadConfig
	Consul         ConsulConfig
	Vault          VaultConfig
	FunctionPrefix string
	ReadTimeout    time.Duration `hcl:"read_timeout"`
	WriteTimeout   time.Duration `hcl:"write_timeout"`
}

type AuthConfig struct {
	Enabled        bool
	Type           string
	CredentialsDir string `hcl:"credentials_dir"`
}

type VaultConfig struct {
	Client  *vaultapi.Client
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
		Scheduling: SchedulingDefaults{
			JobPrefix: "openfaas",
		},
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

	// Merge both decoded config and default config
	err = mergo.Merge(decoded, pc)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (pc *ProviderConfig) LoadCommandLine(args map[string]interface{}) *ProviderConfig {
	pc.ListenPort = args[ListenPort].(int)
	pc.Consul.Address = args[ConsulAddr].(string)
	pc.Consul.TLS.Insecure = args[ConsulTLSSkipVerify].(bool)
	pc.Nomad.Address = args[NomadAddr].(string)
	pc.Nomad.TLS.Insecure = args[NomadTLSSkipVerify].(bool)
	pc.Vault.Address = args[VaultAddr].(string)
	pc.Vault.TLS.Insecure = args[VaultTLSSkipVerify].(bool)
	return pc
}

func (pc *ProviderConfig) MakeVaultClient() error {
	config := vaultapi.DefaultConfig()
	config.ConfigureTLS(&vaultapi.TLSConfig{
		CACert:     pc.Vault.TLS.CAFile,
		ClientCert: pc.Vault.TLS.CertFile,
		ClientKey:  pc.Vault.TLS.KeyFile,
		Insecure:   pc.Vault.TLS.Insecure,
	})
	config.Address = pc.Vault.Address
	client, err := vaultapi.NewClient(config)

	if err != nil {
		return err
	} else {
		pc.Vault.Client = client
		return nil
	}
}

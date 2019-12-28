package types

import "github.com/hashicorp/vault/api"

type VaultConfig struct {
	Address string
	TLS     TLSConfig
	AppRole AppRoleConfig
	Secrets SecretConfig
}

func NewVaultConfig() *VaultConfig {
	vc := &VaultConfig{}
	return vc
}

func (vc *VaultConfig) MakeClient() (Client, error) {
	config := api.DefaultConfig()
	config.ConfigureTLS(&api.TLSConfig{
		CACert:     vc.TLS.CAFile,
		ClientCert: vc.TLS.CertFile,
		ClientKey:  vc.TLS.KeyFile,
		Insecure:   vc.TLS.Insecure,
	})
	config.Address = vc.Address
	if client, err := api.NewClient(config); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

package types

import vaultapi "github.com/hashicorp/vault/api"

type VaultConfig struct {
	Client  *vaultapi.Client
	Address string
	TLS     TLSConfig
	AppRole AppRoleConfig
	Secrets SecretConfig
}

func NewVaultConfig() *VaultConfig {
	nc := &VaultConfig{}
	return nc
}

func (vc *VaultConfig) MakeClient() error {
	config := vaultapi.DefaultConfig()
	config.ConfigureTLS(&vaultapi.TLSConfig{
		CACert:     vc.TLS.CAFile,
		ClientCert: vc.TLS.CertFile,
		ClientKey:  vc.TLS.KeyFile,
		Insecure:   vc.TLS.Insecure,
	})
	config.Address = vc.Address
	client, err := vaultapi.NewClient(config)

	if err != nil {
		return err
	} else {
		vc.Client = client
		return nil
	}
}

package types

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Client  *api.Client
	Address string
	TLS     TLSConfig
	AppRole AppRoleConfig
	Secrets SecretConfig
}

func NewVaultConfig() *VaultConfig {
	vc := &VaultConfig{}
	return vc
}

func (vc *VaultConfig) MakeClient() (*api.Client, error) {
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
		vc.Client = client
		return client, nil
	}
}

func (vc *VaultConfig) Login() (*api.Secret, error) {
	var vaultLogin api.Secret

	lResp, err := vc.DoRequest("POST", fmt.Sprintf("/v1/auth/%s/login", vc.AppRole.Path),
		map[string]interface{}{"role_id": vc.AppRole.RoleID, "secret_id": vc.AppRole.SecretID})

	if err != nil {
		return &vaultLogin, err
	}

	lBody, _ := ioutil.ReadAll(lResp.Body)
	err = json.Unmarshal(lBody, &vaultLogin)
	if err != nil {
		return &vaultLogin, err
	}

	if vaultLogin.Auth != nil && len(vaultLogin.Auth.ClientToken) > 0 {
		vc.Client.SetToken(vaultLogin.Auth.ClientToken)

		r, err := vc.Client.NewRenewer(&api.RenewerInput{
			Secret: &vaultLogin,
		})
		if err != nil {
			return nil, err
		}
		go r.Renew()
	}

	return &vaultLogin, nil
}

func (vc *VaultConfig) DoRequest(method string, path string, body interface{}) (*http.Response, error) {

	if vc.Client == nil {
		_, err := vc.MakeClient()
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{}
	trIgnore := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	createRequest := vc.Client.NewRequest(method, path)
	createRequest.SetJSONBody(body)

	request, _ := createRequest.ToHTTP()
	if vc.TLS.Insecure {
		client.Transport = trIgnore
	}
	return client.Do(request)
}

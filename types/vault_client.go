package types

import "github.com/hashicorp/vault/api"

type Client interface {
	SetAddress(addr string) error
	NewRequest(method, requestPath string) *api.Request
	SetToken(v string)
	NewRenewer(i *api.RenewerInput) (*api.Renewer, error)
}

type Renewer interface {
	Renew()
}

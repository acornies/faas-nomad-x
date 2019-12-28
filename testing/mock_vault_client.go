package testing

import (
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/mock"
)

type MockVaultClient struct {
	mock.Mock
}

func DefaultConfig() *api.Config {
	return &api.Config{}
}

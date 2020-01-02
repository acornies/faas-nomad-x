package testing

import (
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/mock"
)

type MockVaultClient struct {
	mock.Mock
}

func (m *MockVaultClient) SetAddress(addr string) error {

	args := m.Called(addr)

	return args.Error(0)
}

func (m *MockVaultClient) NewRequest(method, requestPath string) *api.Request {

	args := m.Called(method, requestPath)

	var req *api.Request
	if r := args.Get(0); r != nil {
		req = r.(*api.Request)
	}

	return req
}

func (m *MockVaultClient) SetToken(v string) {

	m.Called(v)
}

func (m *MockVaultClient) NewRenewer(i *api.RenewerInput) (*api.Renewer, error) {

	args := m.Called(i)

	var resp *api.Renewer
	if r := args.Get(0); r != nil {
		resp = r.(*api.Renewer)
	}

	return resp, args.Error(1)
}

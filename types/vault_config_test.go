package types

import "testing"

func TestMakeVaultClientWithDefaults(t *testing.T) {
	vc := NewVaultConfig()
	err := vc.MakeClient()
	if err != nil {
		t.Error("Unexpected failure to create Vault client using default configuration", err)
	}
	if vc.Client == nil {
		t.Error("Unexpected nil reference in ProviderConfig.Vault.Client")
	}
}

func TestMakeVaultClientBadConfig(t *testing.T) {
	vc := NewVaultConfig()
	vc.Address = "bad!@#$%^&*()_{address}"
	err := vc.MakeClient()
	if err == nil {
		t.Error("Unexpected success in creation of Vault client")
	}
}

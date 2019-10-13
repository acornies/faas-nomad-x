package types

import "testing"

func TestMakeNomadClientWithDefaults(t *testing.T) {
	nc := NewNomadConfig()
	jobs, err := nc.MakeClient()
	if err != nil {
		t.Error("Unexpected failure to create Nomad client using default configuration", err)
	}
	if jobs == nil {
		t.Error("Unexpected nil reference returning Nomad client")
	}
}

func TestMakeNomadClientBadAddress(t *testing.T) {
	nc := NewNomadConfig()
	nc.Address = "127.0.1.1:4646"
	_, err := nc.MakeClient()
	if err == nil {
		t.Error("Unexpected success in creation of Nomad client")
	}
}

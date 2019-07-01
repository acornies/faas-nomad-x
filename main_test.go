package main

import "testing"

func TestConfigureDefault(t *testing.T) {
	pc, _ := configure("")
	if pc == nil {
		t.Error("Failed to return default configuration to provider")
	}
}

func TestConfigureWithFile(t *testing.T) {
	_, err := configure("default.hcl")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

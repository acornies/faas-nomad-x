package main

import "testing"

func TestConfigureDefault(t *testing.T) {
	configFile := ""
	pc, _ := configure(&configFile)
	if pc == nil {
		t.Error("Failed to return default configuration to provider")
	}
}

func TestConfigureWithFile(t *testing.T) {
	configFile := "default.hcl"
	_, err := configure(&configFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

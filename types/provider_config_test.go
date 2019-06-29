package types

import (
	"path/filepath"
	"testing"
)

func setupTest() *ProviderConfig {
	return NewProviderConfig()
}

func TestLoadFileBadSyntax(t *testing.T) {
	testFile := "bad-config.hcl"
	path, _ := filepath.Abs(filepath.Join("./test-fixtures", testFile))
	pc, err := setupTest().LoadFile(path)

	if err == nil {
		t.Errorf("Bad HCL file expected but was successful: %s", path)
	}

	t.Logf("Bad syntax in .hcl file, fallback to default: %v %v", pc, err)
}

func TestLoadFileBadPath(t *testing.T) {
	testFile := "missing-config.hcl"
	path, _ := filepath.Abs(filepath.Join("./test-fixtures", testFile))
	pc, err := setupTest().LoadFile(path)

	if err == nil {
		t.Errorf("Missing file exists: %s", testFile)
	}
	t.Logf("Missing file, fallback to default config: %v %v", pc, err)
}

func TestLoadFileListenPort(t *testing.T) {
	testFile := "listen-config.hcl"
	path, _ := filepath.Abs(filepath.Join("./test-fixtures", testFile))
	pc, err := setupTest().LoadFile(path)

	if err != nil {
		t.Fatalf("Failed to load test file: %s", testFile)
	}

	t.Logf("test file %s", path)
	if pc.ListenPort != 8082 {
		t.Errorf("Unexpected listen port: %d", pc.ListenPort)
	}
}

func TestLoadFileExampleWithDefaults(t *testing.T) {
	testFile := "example-config.hcl"
	path, err := filepath.Abs(filepath.Join("./test-fixtures", testFile))
	if err != nil {
		t.Fatalf("Failed to load test file: %s", testFile)
	}
	t.Logf("test file %s", path)

	pc, err := setupTest().LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load file: %s", path)
	}

	if pc.ListenPort != 8081 {
		t.Errorf("Unexpected listen port: %d", pc.ListenPort)
	}
	if pc.LogLevel != "DEBUG" {
		t.Errorf("Unexpected log_level port: %s", pc.LogLevel)
	}
	if pc.Nomad.Address != "127.0.0.1:4646" {
		t.Errorf("Unexpected Nomad address: %s", pc.Nomad.Address)
	}
	if pc.Nomad.Driver != "exec" {
		t.Errorf("Unexpected Nomad driver: %s", pc.Nomad.Driver)
	}
	if pc.Nomad.ACLToken != "abcdefg" {
		t.Errorf("Unexpected Nomad ACL: %s", pc.Nomad.ACLToken)
	}
	if pc.Nomad.Region != "east-us" {
		t.Errorf("Unexpected Nomad region: %s", pc.Nomad.Region)
	}
}

func TestLoadCommandLine(t *testing.T) {
	pc := setupTest().LoadCommandLine(8080, "127.0.1.1:8500", "127.0.1.1:4646", "127.0.1.1:8200")

	if pc.ListenPort != 8080 {
		t.Errorf("Unexpected listen port from cli: %d", pc.ListenPort)
	}
	if pc.Nomad.Address != "127.0.1.1:4646" {
		t.Errorf("Unexpected Nomad address from cli: %s", pc.Nomad.Address)
	}
	if pc.Consul.Address != "127.0.1.1:8500" {
		t.Errorf("Unexpected Consul address from cli: %s", pc.Consul.Address)
	}
	if pc.Vault.Address != "127.0.1.1:8200" {
		t.Errorf("Unexpected Vault address from cli: %s", pc.Vault.Address)
	}
}

func TestDefault(t *testing.T) {
	pc := setupTest().Default()
	if pc.Nomad.Address != "127.0.0.1:4646" {
		t.Errorf("Unexpected Nomad default address: %s", pc.Nomad.Address)
	}
	if pc.Consul.Address != "127.0.0.1:8500" {
		t.Errorf("Unexpected Consul default address: %s", pc.Consul.Address)
	}
	if pc.Vault.Address != "127.0.0.1:8200" {
		t.Errorf("Unexpected Vault default address: %s", pc.Vault.Address)
	}
}

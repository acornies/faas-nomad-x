package types

import "github.com/hashicorp/nomad/api"

type NomadConfig struct {
	Address    string
	ACLToken   string `hcl:"acl_token"`
	TLS        TLSConfig
	Region     string
	Namespace  string
	Driver     string
	Datacenter string
	Scheduling SchedulingDefaults
}

type SchedulingDefaults struct {
	JobType         string `hcl:"job_type"`
	JobPrefix       string `hcl:"job_prefix"`
	Replicas        int
	Memory          int
	CPU             int
	RestartAttempts int    `hcl:"restart_attempts"`
	RestartMode     string `hcl:"restart_mode"`
	RestartDelay    string `hcl:"restart_delay"`
	Priority        int
	DiskSize        int
	NetworkingMode  string `hcl:"network_mode"`
}

func NewNomadConfig() *NomadConfig {
	nc := &NomadConfig{}
	return nc
}

func (nc *NomadConfig) MakeClient() (Jobs, error) {
	client, err := api.NewClient(&api.Config{
		Address:  nc.Address,
		Region:   nc.Region,
		SecretID: nc.ACLToken,
		TLSConfig: &api.TLSConfig{
			CACert:     nc.TLS.CAFile,
			ClientCert: nc.TLS.CertFile,
			ClientKey:  nc.TLS.KeyFile,
			Insecure:   nc.TLS.Insecure,
		},
	})
	if err != nil {
		return nil, err
	} else {
		if len(nc.Datacenter) <= 0 {
			nc.Datacenter, _ = client.Agent().Datacenter()
		}
		return client.Jobs(), nil
	}
}

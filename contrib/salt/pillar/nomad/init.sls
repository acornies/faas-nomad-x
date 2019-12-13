nomad:
  version: '0.10.2'
  config_dir: '/etc/nomad.d'
  bin_dir: '/usr/local/bin'
  service_hash: 1aea4ba5283cd79264da5c7e3214049f86f77af4
  config:
    datacenter: dc1
    data_dir: /var/lib/nomad
    log_level: DEBUG
    server:
      enabled: true
      bootstrap_expect: 1
      encrypt: "AaABbB+CcCdDdEeeFFfggG=="
    addresses:
      http: {{ grains['ip_interfaces']['enp0s8'][0] }}
      rpc: {{ grains['ip_interfaces']['enp0s8'][0] }}
      serf: {{ grains['ip_interfaces']['enp0s8'][0] }}
    client:
      network_interface: enp0s8
      enabled: true
      meta:
        service_host: "true"
        faas_host: "true"
    consul:
      address: "{{ grains['ip_interfaces']['enp0s8'][0] }}:8500"
      server_service_name: "nomad"
      client_service_name: "nomad-client"
      auto_advertise: true
      server_auto_join: true
      client_auto_join: true
    vault:
      enabled: true
      address: "http://127.0.0.1:8200"
      token: vagrant
      create_from_role: nomad-cluster
  datacenters:
    - dc1
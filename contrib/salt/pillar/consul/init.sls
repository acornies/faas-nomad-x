consul:
  version: 1.6.0
  # Need to run as root to bind to 53 for dns
  user: root
  # Start Consul agent service and enable it at boot time
  service: True
  config:
    server: True
    data_dir: "/data"
    log_level: DEBUG
    datacenter: dc1
    encrypt: "RIxqpNlOXqtr/j4BgvIMEw=="
    bootstrap_expect: 1
    disable_update_check: True
    enable_syslog: True
    advertise_addr: {{ grains['ip_interfaces']['enp0s8'][0] }}
    ports:
      dns: 53
      http: 8500
      grpc: 8502
    addresses:
      http: {{ grains['ip_interfaces']['enp0s8'][0] }}
      dns: {{ grains['ip_interfaces']['enp0s8'][0] }}
      grpc: {{ grains['ip_interfaces']['enp0s8'][0] }}
    connect:
      enabled: True
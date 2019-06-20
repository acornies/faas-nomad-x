nomad:
  version: '0.9.3'
  config_dir: '/etc/nomad.d'
  bin_dir: '/usr/local/bin'
  service_hash: 9ac47b2dd92b232dc3f55b03ba0e756ed0e7fd95
  config:
    datacenter: dc1
    # tls:
    #   http: True
    #   ca_file: /home/vagrant/placeholder-ca.crt
    #   cert_file: /home/vagrant/placeholder.crt
    #   key_file: /home/vagrant/placeholder.key
    data_dir: /var/lib/nomad
    log_level: DEBUG
    server:
      enabled: true
      bootstrap_expect: 1
      encrypt: "AaABbB+CcCdDdEeeFFfggG=="
    addresses:
      {% if grains['provider'] == 'virtualbox' %}
      http: {{ grains['ip_interfaces']['enp0s3'][0] }}
      rpc: {{ grains['ip_interfaces']['enp0s3'][0] }}
      serf: {{ grains['ip_interfaces']['enp0s3'][0] }}
      {% elif grains['provider'] == 'vmware' %}
      http: {{ grains['ip_interfaces']['eth0'][0] }}
      rpc: {{ grains['ip_interfaces']['eth0'][0] }}
      serf: {{ grains['ip_interfaces']['eth0'][0] }}
      {% elif grains['provider'] == 'libvirt' %}
      http: {{ grains['ip_interfaces']['eth0'][0] }}
      rpc: {{ grains['ip_interfaces']['eth0'][0] }}
      serf: {{ grains['ip_interfaces']['eth0'][0] }}
      {% endif %}
    client:
      {% if grains['provider'] == 'virtualbox' %}
      network_interface: enp0s3
      {% elif grains['provider'] == 'vmware' %}
      network_interface: eth0
      {% elif grains['provider'] == 'libvirt' %}
      network_interface: eth0
      {% endif %}
      enabled: true
      meta:
        service_host: "true"
        faas_host: "true"
    consul:
      address: "127.0.0.1:8500"
      server_service_name: "nomad"
      client_service_name: "nomad-client"
      auto_advertise: true
      server_auto_join: true
      client_auto_join: true
    vault:
      enabled: true
      address: "http://127.0.0.1:8200"
      token: vagrant
  datacenters:
    - dc1
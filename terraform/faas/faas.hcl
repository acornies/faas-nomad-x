job "faas" {

  datacenters = ["dc1"]

  type = "system"

  group "faas-svc" {

    task "faas-gateway" {
      driver = "docker"
      template {
        env = true
        destination  = "secrets/gateway.env"
        data = <<EOH
functions_provider_url="http://{{ env "NOMAD_IP_http" }}:8081/"
faas_prometheus_host="{{ env "NOMAD_IP_http" }}"
faas_prometheus_port="9090"
{{ range service "nats" }}
faas_nats_address="{{ .Address }}"
faas_nats_port={{ .Port }}{{ end }}
read_timeout="5m5s" # Maximum time to read HTTP request
write_timeout="5m5s" # Maximum time to write HTTP response
upstream_timeout="5m" # Maximum duration of upstream function call - should be more than read_timeout and write_timeout
dnsrr="false" # Temporarily use dnsrr in place of VIP while issue persists on PWD
direct_functions="false" # Functions are invoked directly over the overlay network
direct_functions_suffix=""
basic_auth="true"
secret_mount_path="/secrets/"
scale_from_zero="false" # Enable if you want functions to scale from 0/0 to min replica count upon invoke
max_idle_conns="1024"
max_idle_conns_per_host="1024"
auth_proxy_url="http://{{ env "NOMAD_IP_http" }}:8083/validate"
auth_proxy_pass_body="false"
EOH
      }

      config {
        image = "openfaas/gateway:0.17.0"
        port_map {
          http = 8080
          metrics = 8082
        }
      }

      vault {
        policies = ["openfaas"]
      }
      // basic auth from vault example
      // update -enable_basic_auth=true
      // uncomment below if you have a Vault instance connected to Nomad
      template {
        destination   = "secrets/basic-auth-user"
        data = <<EOH
{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.username }}{{ end }}
EOH
      }
      template {
        destination   = "secrets/basic-auth-password"
        data = <<EOH
{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.password }}{{ end }}
EOH
      }

      resources {
        network {
          port "http" {
            static = 8080
          }
          port "metrics" {
            static = 8082
          }
        }
      }

      service {
        port = "http"
        name = "gateway"
        tags = ["faas"]
      }

      service {
        port = "metrics"
        name = "faas-metrics"
        tags = ["faas"]
      }
    }

    task "basic-auth-plugin" {
      driver = "docker"
      template {
        env = true
        destination   = "secrets/gateway.env"

        data = <<EOH
secret_mount_path="/secrets/"
user_filename="basic-auth-user"
pass_filename="basic-auth-password"
EOH
      }

      config {
        image = "openfaas/basic-auth-plugin:0.17.0"
        port_map {
          http = 8080
        }
      }

      vault {
        policies = ["openfaas"]
      }
      // basic auth from vault example
      // update -enable_basic_auth=true
      // uncomment below if you have a Vault instance connected to Nomad
      template {
        destination   = "secrets/basic-auth-user"
        data = <<EOH
{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.username }}{{ end }}
EOH
      }
      template {
        destination   = "secrets/basic-auth-password"
        data = <<EOH
{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.password }}{{ end }}
EOH
      }

      resources {
        memory = 50
        network {
          port "http" {
            static = 8083
          }
        }
      }
    }

    task "statsd" {
      driver = "docker"

      config {
        image = "prom/statsd-exporter:v0.12.2"

        args = [
          "--log.level=debug",
        ]
      }

      resources {
        network {

          port "http" {
            static = 9102
          }

          port "statsd" {
            static = 9125
          }
        }
      }

      service {
        port = "http"
        name = "statsd"
        tags = ["faas"]

        check {
          type     = "http"
          port     = "http"
          interval = "10s"
          timeout  = "2s"
          path     = "/"
        }
      }
    }
  }

  group "faas-nats" {

    task "nats" {
      driver = "docker"
      
      config {
        image = "nats-streaming:0.11.2"

        args = [
          "-store", "file", "-dir", "/tmp/nats",
          "-m", "8222",
          "-cid","faas-cluster",
        ]

        port_map {
          client = 4222,
          monitoring = 8222
          routing = 6222
        }
      }

      resources {
        memory = 100
        network {
          port "client" {
            static = 4222
          }

          port "monitoring" {
            static = 8222
          }

          port "routing" {
            static = 6222
          }
        }
      }

      service {
        port = "client"
        name = "nats"
        tags = ["faas"]

        check {
           type     = "http"
           port     = "monitoring"
           path     = "/connz"
           interval = "5s"
           timeout  = "2s"
        }
      }
    }
  }

  group "faas-monitoring" {

    task "prometheus" {
      driver = "docker"
      
      config {
        image = "prom/prometheus:v2.7.1"

        args = [
          "--config.file=/local/prometheus.yml"
        ]

        port_map {
          http = 9090
        }
      }

      artifact {
			  source      = "https://raw.githubusercontent.com/acornies/THUG-aug27-2019/master/contrib/prometheus/prometheus.yml"
			  destination = "local/prometheus.yml.tpl"
				mode        = "file"
			}

      template {
        source        = "local/prometheus.yml.tpl"
        destination   = "local/prometheus.yml"
        change_mode   = "noop"
      }
			
			artifact {
			  source      = "https://raw.githubusercontent.com/acornies/THUG-aug27-2019/master/contrib/prometheus/alert.rules.yml"
			  destination = "local/alert.rules.yml"
				mode        = "file"
			}

      resources {
        network {
          port "http" {
            static = 9090
          }
        }
      }

      service {
        port = "http"
        name = "prometheus"
        tags = ["faas"]

        check {
          type     = "http"
          port     = "http"
          interval = "10s"
          timeout  = "2s"
          path     = "/graph"
        }
      }
    }

    task "alertmanager" {
      driver = "docker"

			artifact {
			  source      = "https://raw.githubusercontent.com/acornies/THUG-aug27-2019/master/contrib/prometheus/alertmanager.yml"
			  destination = "local/alertmanager.yml.tpl"
				mode        = "file"
			}

      template {
        source        = "local/alertmanager.yml.tpl"
        destination   = "local/alertmanager.yml"
        change_mode   = "noop"
      }

      config {
        image = "prom/alertmanager:v0.16.1"

        port_map {
          http = 9093
        }

        args = [
          "--config.file=/local/alertmanager.yml"
        ]
      }

      vault {
        policies = ["default", "openfaas"]
      }
      // basic auth from vault example
      // update -enable_basic_auth=true
      // uncomment below if you have a Vault instance connected to Nomad
      template {
        destination   = "secrets/basic-auth-user"
        data = <<EOH
{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.username }}{{ end }}
EOH
      }
      template {
        destination   = "secrets/basic-auth-password"
        data = <<EOH
{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.password }}{{ end }}
EOH
      }

      resources {
        network {
          port "http" {
            static = 9093
          }
        }
      }

      service {
        port = "http"
        name = "alertmanager"
        tags = ["faas"]
      }
    }
  }
}
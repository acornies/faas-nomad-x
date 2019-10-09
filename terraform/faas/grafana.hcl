job "grafana" {

  type = "service"

  datacenters = ["dc1"]

  group "grafana-instances" {
    count = 1

    vault {
      policies = ["default", "openfaas"]
    }

    task "grafana-svc" {        
      driver = "docker"
      config {
        image = "grafana/grafana:6.3.2"
        port_map = {         
          http = 3000
        }
        volumes = ["local/provisioning:/provisioning"]
      }

      env {            
        GF_SERVER_PROTOCOL="http"
        GF_SERVER_ROOT_URL="http://$${NOMAD_IP_http}:3000"
        GF_PATHS_PROVISIONING="/provisioning"
      }

      template {
        data = <<EOH
# SECURITY
GF_SECURITY_ADMIN_USER="{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.username }}{{ end }}"
GF_SECURITY_ADMIN_PASSWORD="{{ with secret "secret/openfaas/auth/credentials" }}{{ .Data.password }}{{ end }}"
EOH
        env = true
        destination = "secrets/grafana.env"
        change_mode = "noop"
        splay = "1m"
      }
      
      template {
        data = <<EOH
apiVersion: 1

datasources:
  - name: prometheus
    type: prometheus
    access: proxy
    url: http://{{ env "NOMAD_IP_http" }}:9090
    isDefault: true
EOH

        destination = "local/provisioning/datasources/datasources.yaml"
        change_mode = "noop"      
        splay = "1m"
      }
      
      service {
        name = "grafana"
        tags = ["grafana"]
        port = "http"
        check {
          type = "tcp"
          port     = "http"
          interval = "10s"
          timeout  = "5s"
        }
      }
      resources {
        network {     
          port "http" {
            static = 3000
          }
        }
      }
    }
  }
}
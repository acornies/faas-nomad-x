job "faas-kafka-connector" {

  datacenters = ["dc1"]

  type = "service"

  group "faas-kafka" {

    task "faas-kafka-connector" {
      driver = "docker"

      config {
        image = "openfaas/kafka-connector:0.4.0"
      }

      vault {
        policies = ["default", "openfaas"]
      }

      env {            
        gateway_url="http://$${attr.unique.network.ip-address}:8080"
        topics="stripe-webhook-charge-dispute-created"
        broker_host="$${attr.unique.network.ip-address}:29092"
        basic_auth="true"
        secret_mount_path="/secrets/"
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
          port "http" {}
        }
      }
    }
  }
}
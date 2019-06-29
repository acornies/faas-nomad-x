listen_port = 8081
log_level = "INFO"
health_enabled = true
auth_enabled = false
// credentials_dir = ""

// Required for function schdeuling
nomad {
  driver = "docker"
  address = "127.0.0.1:4646"
  acl_token = ""
  region = "global"
  tls {
    ca_file = ""
    cert_file = ""
    key_file = ""
    insecure = false
  }
}

// Required for faas-cli secret API
vault {
  address = "127.0.0.1:8200"
  app_role {
    role_id = ""
    secret_id = ""
  }
  secrets {
    key_prefix = "kv/openfaas"
    // change to 1 to use kv backend v1
    kv_version = 2
    policy = "openfaas"
  }
  tls {
    ca_file = ""
    cert_file = ""
    key_file = ""
    insecure = false
  }
}

// Required for function service discovery
consul {
  address = "127.0.0.1:8500"
  acl_token = "placeholder"
  dns_enabled = false
  tls {
    ca_file = ""
    cert_file = ""
    key_file = ""
    insecure = false
  }
}
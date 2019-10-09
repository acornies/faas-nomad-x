resource "vault_policy" "openfaas" {
  name = "openfaas"

  policy = <<EOT
path "secret/openfaas/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
EOT

}

resource "vault_approle_auth_backend_role" "openfaas" {
  backend        = "approle/"
  role_name      = "openfaas"
  token_policies = [vault_policy.openfaas.name]
  // 5 mins
  token_period = 300
}

resource "vault_approle_auth_backend_role_secret_id" "faas_nomad" {
  backend   = "approle/"
  role_name = vault_approle_auth_backend_role.openfaas.role_name
}


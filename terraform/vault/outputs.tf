output "role_id" {
  value = [vault_approle_auth_backend_role.openfaas.role_id]
}

output "secret_id" {
  sensitive = true
  value     = [vault_approle_auth_backend_role_secret_id.faas_nomad.secret_id]
}


module "vault" {
  source = "./vault"
}

module "faas" {
  source               = "./faas"
  vault_approle_id     = module.vault.role_id[0]
  vault_approle_secret = module.vault.secret_id[0]
}

// module "kafka" {
//   source = "./kafka"
// }

// module "faas_kafka_connector" {
//   source = "./kafka_connector"
// }

provider "nomad" {
  address = "http://192.168.50.2:4646"
}

provider "vault" {
  address = "http://192.168.50.2:8200"
  token   = var.vault_token
}


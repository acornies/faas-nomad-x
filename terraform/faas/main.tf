data "template_file" "faas" {
  template = file("${path.module}/faas.hcl")
  vars = {
    vault_approle_id     = var.vault_approle_id
    vault_approle_secret = var.vault_approle_secret
  }
}

resource "nomad_job" "faas" {
  jobspec = data.template_file.faas.rendered
}

data "template_file" "grafana" {
  template = file("${path.module}/grafana.hcl")
}

resource "nomad_job" "grafana" {
  jobspec = data.template_file.grafana.rendered
}


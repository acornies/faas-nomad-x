data "template_file" "kafka" {
  template = file("${path.module}/kafka.hcl")
}

resource "nomad_job" "kafka" {
  jobspec = data.template_file.kafka.rendered
}


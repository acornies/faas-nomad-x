data "template_file" "kafka_connector" {
  template = file("${path.module}/connector.hcl")
}

resource "nomad_job" "kafka_connector" {
  jobspec = data.template_file.kafka_connector.rendered
}


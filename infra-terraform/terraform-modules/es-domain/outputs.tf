output "es_point" {
  value = aws_elasticsearch_domain.es.endpoint
}
output "kibana_endpoint" {
  value = aws_elasticsearch_domain.es.kibana_endpoint
}
output "masterusername" {
  value = var.masterusername
}
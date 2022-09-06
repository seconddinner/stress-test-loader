output "vpc_id" {
  value = module.network.aws_vpc_id
}
output "aws_subnets" {
  value = module.network.aws_subnets
}
output "subnet_ids" {
  value = tolist(data.aws_subnet_ids.self.ids)
}
output "es_point" {
  value = aws_elasticsearch_domain.es.endpoint
}
output "kibana_endpoint" {
  value = aws_elasticsearch_domain.es.kibana_endpoint
}
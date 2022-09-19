output "vpc_id" {
  value = module.network.aws_vpc_id
}
#output "aws_subnets" {
#  value = module.network.aws_subnets
#}
output "subnet_ids" {
  value = data.aws_subnets.self.ids
}
output "es_endpoint" {
  value = var.create_esdomain ? module.es-domain[0].es_endpoint : data.aws_elasticsearch_domain.es[0].endpoint
}
output "es_kibana_endpoint" {
  value = var.create_esdomain ? module.es-domain[0].kibana_endpoint : data.aws_elasticsearch_domain.es[0].kibana_endpoint
}
output "es_masterusername" {
  value = var.masterusername
}
output "stressclient_asg_name" {
  value = module.autoscale.stressclient_asg_name
}
output "stressclient_asg_desired_capacity" {
  value = module.autoscale.stressclient_asg_desired_capacity
}

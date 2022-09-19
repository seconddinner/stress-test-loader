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
  value = module.es-domain.es_endpoint
}
output "es_kibana_endpoint" {
  value = module.es-domain.kibana_endpoint
}
output "es_masterusername" {
  value = module.es-domain.masterusername
}
output "stressclient_asg_name" {
  value = module.autoscale.stressclient_asg_name
}
output "stressclient_asg_desired_capacity" {
  value = module.autoscale.stressclient_asg_desired_capacity
}

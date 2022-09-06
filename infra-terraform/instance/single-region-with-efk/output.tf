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
  value = module.es-domain.es_point
}
output "kibana_endpoint" {
  value = module.es-domain.kibana_endpoint
}
output "masterusername" {
  value = module.es-domain.masterusername
}
output "stressclient_asg_name" {
  value = module.autoscale.stressclient_asg_name
}
output "stressclient_asg_desired_capacity" {
  value = module.autoscale.stressclient_asg_desired_capacity
}

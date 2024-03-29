output "vpc_id" {
  value = module.network.aws_vpc_id
}
#output "aws_subnets" {
#  value = module.network.aws_subnets
#}
output "subnet_ids" {
  value = data.aws_subnets.self.ids
}
output "stressclient_asg_name" {
  value = module.autoscale.stressclient_asg_name
}
output "stressclient_asg_desired_capacity" {
  value = module.autoscale.stressclient_asg_desired_capacity
}
output "stressclient_ami" {
  value = data.aws_ami.stl.image_id
}

# output "public-ips" {
#   value = "${data.aws_instances.workers.public_ips}"
# }

output "stressclient_asg_name" {
  value = aws_autoscaling_group.stress_test_loader.name
}
output "stressclient_asg_desired_capacity" {
  value = aws_autoscaling_group.stress_test_loader.desired_capacity
}
output "aws_subnet_list" {
  value = ["${aws_subnet.main.*.id}"]
  # value = aws_subnet.main.id
}

output "aws_vpc_id" {
  value = aws_vpc.main.id
}

output "aws_subnets" {
  value = aws_subnet.main
}

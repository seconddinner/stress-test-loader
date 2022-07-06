output "nlb_dns" {
  value = aws_lb.nginx
}

output "dns" {
  value = "${aws_route53_record.domain.name}.${var.domain_root}"
}

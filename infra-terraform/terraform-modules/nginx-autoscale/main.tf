

resource "aws_key_pair" "sd_stresstest" {
  key_name   = "sd_stresstest_${var.environment}_${var.domain_rand}"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC03A2vpgWrbqFechOFlT+thhD40XlxH9Dc0WvxtioJUm4KASFYc0hQ+xkJGdciiCWyBvgqRU9x2bE+8r5S4EYpnXg+hPCEiyOfXT204HeUJjr6zZwn7NOvlyBED/pCn4sP2JcX2pJNDKlJTak3fZGCqndjy7ue8AZj/hZf1sbIJqojjRXEdyT/eE/CJkozWaI30AXDsoFUFjtxZzuMPD3fdj7Vpxz6kK0XCWnF0ObbjsZcfUUFphSZiI9so5t7wBJ2XsL0IB6/VjeGBuGS0VRg+QNnCziQOukzFiq67YFdDMDsNWnh04028qCnZfTtFEM3bKw18RvgfpB0QpL5ZJu0aTtfTB2zeg3FoE/N3qfemuJsewqYtCYV6LlH2RuNeZcfj1ZCyWazzF/XCkwyMnTprT+TdXIMGj1gmbnxWnjvkNiYCpEIeo8CrMLiZtj/03pin26+hHb7a8CVFXOTZZktf2UtN7sOwSfIZysvJbiEKSZ58Khv1Pn0AJjg3dexueE="
}


resource "aws_launch_configuration" "nginx" {
  name_prefix   = "nginx-${var.environment}"
  image_id      = var.aws_ami_id
  instance_type = var.instance_type
  # associate_public_ip_address = false
  key_name        = aws_key_pair.sd_stresstest.key_name
  security_groups = ["${aws_security_group.instance.id}"]

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 30
    delete_on_termination = true
  }

  #user_data = data.template_file.user_data.rendered
  user_data = var.user_data

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "instance" {
  name   = "nginx-instance-${var.environment}"
  vpc_id = var.vpc_id

  # ping
  ingress {
    from_port   = 0
    to_port     = 8
    protocol    = "icmp"
    cidr_blocks = var.ssh_allowed_cidr
  }

  # ssh
  ingress {
    from_port   = "22"
    to_port     = "22"
    protocol    = "TCP"
    cidr_blocks = var.ssh_allowed_cidr
  }

  # nginx
  ingress {
    from_port   = var.nginx_port
    to_port     = var.nginx_port
    protocol    = "TCP"
    cidr_blocks = var.nginx_allowed_cidr
  }

  # prometheus-node-exporter
  ingress {
    from_port   = "9100"
    to_port     = "9100"
    protocol    = "TCP"
    cidr_blocks = var.prometheus_node_allowed_cidr
  }

  # vpc cidr_block to allow lb check instance
  ingress {
    from_port   = var.nginx_port
    to_port     = var.nginx_port
    protocol    = "TCP"
    cidr_blocks = [var.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["${var.egress_allowed_cidr}"]
  }
}

resource "aws_lb" "nginx" {
  internal                         = false
  load_balancer_type               = "network"
  enable_cross_zone_load_balancing = true
  subnets                          = var.aws_subnets.*.id
}

resource "aws_lb_target_group" "tcp" {
  protocol             = "TCP"
  port                 = var.nginx_port
  vpc_id               = var.vpc_id
  deregistration_delay = 60
  health_check {
    protocol            = "TCP"
    port                = var.nginx_port
    interval            = 10
    unhealthy_threshold = 2
    healthy_threshold   = 2
  }
}

resource "aws_lb_listener" "tcp" {
  load_balancer_arn = aws_lb.nginx.arn
  protocol          = "TCP"
  port              = var.nginx_port

  default_action {
    target_group_arn = aws_lb_target_group.tcp.arn
    type             = "forward"
  }
}

resource "aws_autoscaling_group" "nginx" {
  launch_configuration = aws_launch_configuration.nginx.name
  vpc_zone_identifier  = var.aws_subnets.*.id

  health_check_type = "ELB"

  target_group_arns = [
    "${aws_lb_target_group.tcp.arn}",
  ]

  lifecycle {
    create_before_destroy = true
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      min_healthy_percentage = 50
      instance_warmup        = 300
    }
    triggers = ["tag"]
  }

  min_size         = var.min_size
  max_size         = var.max_size
  desired_capacity = var.desired_capacity

  enabled_metrics = tolist(["GroupMinSize", "GroupMaxSize", "GroupDesiredCapacity", "GroupInServiceInstances", "GroupPendingInstances", "GroupStandbyInstances", "GroupTerminatingInstances", "GroupTotalInstances"])

}

# auto scale up policy
resource "aws_autoscaling_policy" "nginx_up" {
  name                   = "nginx-${var.environment}-${var.PNS_version}-up"
  scaling_adjustment     = var.up_scaling_adjustment
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.nginx.name
}

# auto scale down policy
resource "aws_autoscaling_policy" "nginx_down" {
  name                   = "nginx-${var.environment}-${var.PNS_version}-down"
  scaling_adjustment     = var.down_scaling_adjustment
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.nginx.name
}

resource "aws_cloudwatch_metric_alarm" "nginx_up" {
  alarm_name          = "nginx-${var.environment}-${var.PNS_version}-up"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "120"
  statistic           = "Average"

  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.nginx.name
  }
  threshold         = "120"
  alarm_description = "Check whether EC2 instance CPU utilisation is over 80% on average"
  alarm_actions     = ["${aws_autoscaling_policy.nginx_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "nginx_down" {
  alarm_name          = "nginx-${var.environment}-${var.PNS_version}-down"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "120"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "120"
  statistic           = "Average"
  threshold           = "20"
  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.nginx.name
  }

  alarm_description = "Check whether EC2 instance CPU utilisation is under 20% on average"
  alarm_actions     = ["${aws_autoscaling_policy.nginx_down.arn}"]
}


data "aws_route53_zone" "domain_root" {
  name         = var.domain_root
  private_zone = false
}

resource "aws_route53_record" "aws_domain" {
  zone_id = data.aws_route53_zone.domain_root.zone_id
  name    = "target.${var.dns_name}.${var.environment}.${data.aws_route53_zone.domain_root.name}"
  type    = "A"
  alias {
    name                   = aws_lb.nginx.dns_name
    zone_id                = aws_lb.nginx.zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "domain" {
  zone_id = data.aws_route53_zone.domain_root.zone_id
  name    = "target.${var.environment}"
  type    = "CNAME"
  ttl     = "5"
  records = [
    aws_route53_record.aws_domain.name
  ]
}

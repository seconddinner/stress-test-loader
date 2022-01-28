

resource "aws_key_pair" "sd_stresstest" {
  key_name   = "sd_stresstest_${var.environment}_${var.domain_rand}"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILbYzCoHIpkhqJIdJCer668wZHW670XucHd22HLrRGWM"
}


resource "aws_launch_configuration" "stress_test_loader" {
  name_prefix                 = "stress_test_loader-${var.environment}"
  image_id                    = var.aws_ami_id
  instance_type               = var.instance_type
  iam_instance_profile = "${aws_iam_instance_profile.stress_test_client_read_profile.name}"
  associate_public_ip_address = true
  key_name                    = aws_key_pair.sd_stresstest.key_name
  security_groups             = ["${aws_security_group.instance.id}"]

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
  name   = "stress_test_loader-instance-${var.environment}"
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

  # stress_test_loader
  ingress {
    from_port   = var.stress_test_loader_port
    to_port     = var.stress_test_loader_port
    protocol    = "TCP"
    cidr_blocks = var.stress_test_loader_allowed_cidr
  }

  # prometheus-node-exporter
  ingress {
    from_port   = "9100"
    to_port     = "9100"
    protocol    = "TCP"
    cidr_blocks = var.prometheus_node_allowed_cidr
  }

  # prometheus-stress_test_loader-exporter
  ingress {
    from_port   = "9301"
    to_port     = "9301"
    protocol    = "TCP"
    cidr_blocks = var.prometheus_node_allowed_cidr
  }

  # vpc cidr_block to allow lb check instance
  ingress {
    from_port   = var.stress_test_loader_port
    to_port     = var.stress_test_loader_port
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


resource "aws_autoscaling_group" "stress_test_loader" {
  launch_configuration = aws_launch_configuration.stress_test_loader.name
  vpc_zone_identifier  = var.aws_subnets.*.id

  health_check_type = "EC2"


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
resource "aws_autoscaling_policy" "stress_test_loader_up" {
  name                   = "stress_test_loader-${var.environment}-${var.PNS_version}-up"
  scaling_adjustment     = var.up_scaling_adjustment
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.stress_test_loader.name
}

# auto scale down policy
resource "aws_autoscaling_policy" "stress_test_loader_down" {
  name                   = "stress_test_loader-${var.environment}-${var.PNS_version}-down"
  scaling_adjustment     = var.down_scaling_adjustment
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.stress_test_loader.name
}

resource "aws_cloudwatch_metric_alarm" "stress_test_loader_up" {
  alarm_name          = "stress_test_loader-${var.environment}-${var.PNS_version}-up"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "120"
  statistic           = "Average"

  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.stress_test_loader.name
  }
  threshold         = "120"
  alarm_description = "Check whether EC2 instance CPU utilisation is over 80% on average"
  alarm_actions     = ["${aws_autoscaling_policy.stress_test_loader_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "stress_test_loader_down" {
  alarm_name          = "stress_test_loader-${var.environment}-${var.PNS_version}-down"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "120"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "120"
  statistic           = "Average"
  threshold           = "20"
  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.stress_test_loader.name
  }

  alarm_description = "Check whether EC2 instance CPU utilisation is under 20% on average"
  alarm_actions     = ["${aws_autoscaling_policy.stress_test_loader_down.arn}"]
}


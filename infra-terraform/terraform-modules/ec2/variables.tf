variable "asg_name" {
  description = "auto scale group name"
  default     = "auto_scale_group"
}

variable "asg_min" {
  description = ""
}

variable "asg_max" {
  description = ""
}

variable "asg_desired" {
  description = ""
}

variable "aws_subnet" {
  type        = any
  description = ""
}

variable "instance_type" {
  description = ""
}

variable "key_name" {
  description = ""
}

variable "aws_region" {
  description = ""
}

variable "aws_vpc_id" {
  description = ""
}

variable "admin_cidr_ingress" {
  description = ""
}

# variable "aws_security_group_lb_sg_id" {
#   description = ""
# }

# variable "aws_iam_instance_profile_app_name" {
#   description = ""
# }

variable "deployname" {
  description = ""
}

variable "instance_count" {
  default = 1
}

variable "stress_test_loader_instance_type" {
  default = "t2.micro"
}

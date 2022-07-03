# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-east-2"
}

variable "ntw_cidr_block" {
  description = "cider block"
  default     = "10.10.0.0/16"
}

variable "az_count" {
  description = "Number of AZs to cover in a given AWS region"
  default     = 2
}

variable "asg_min" {
  description = ""
  default     = 2
}

variable "asg_max" {
  description = ""
  default     = 500
}

variable "asg_desired" {
  description = ""
  default     = 2
}



variable "big_asg_min" {
  description = ""
  default     = 2
}

variable "big_asg_max" {
  description = ""
  default     = 1000
}

variable "big_asg_desired" {
  description = ""
  default     = 2
}


variable "instance_type" {
  default = "t4g.nano"
}

variable "source_ami_id" {
  default = ""
}

variable "environment" {
  default = "production"
}

variable "stress_test_loader_port" {
  default = 9005
}

variable "PNS_version" {
  type    = string
  default = ""
}

variable "stress_test_loader_allowed_cidr" {
  type = list(string)
}

# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-west-2"
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
  default     = 120
}

variable "asg_desired" {
  description = ""
  default     = 100
}

variable "instance_type" {
  default = "t4g.nano"
}

variable "key_name" {
  default = "sd_stresstest"
}

variable "ami_name" {
  default = ""
}

variable "environment" {
  default = "qa"
}

variable "stress_test_loader_port" {
  default = 9005
}

variable "PNS_version" {
  type    = string
  default = ""
}

variable "stress_test_loader_allowed_cidr" {
  type    = list(string)
  default = ["47.154.122.36/32"]
}

variable "owner_id" {
  default = "275684003332"
}

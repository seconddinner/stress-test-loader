# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."

  type = string
  # default = "ap-southeast-1"
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
  type = list(string)
  default = ["10.0.0.0/8"]
}

variable "owner_id" {
  type = string
  default = ""
}

variable "public_key" {
  type = string
  default = ""
}


variable "esdomain" {
  type = string
  default = "stresstest"
}

variable "masterpassword" {
  type = string
}
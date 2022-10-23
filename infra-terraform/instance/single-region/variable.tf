# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."
  type        = string
  default     = "us-west-2"
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
  default     = 1
}

variable "asg_max" {
  description = ""
  default     = 5
}

variable "asg_desired" {
  description = ""
  default     = 1
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

variable "clientarch" {
  description = "the architecture of the stress-client, valid values are arm or x86_64"
  type        = string
  default     = "arm64"
  validation {
    condition     = contains(["arm64", "x86_64"], var.clientarch)
    error_message = "Valid values for var 'client' are ('arm64', 'x86_64')."
  }
}

variable "environment" {
  default = "gitaction"
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

variable "owner_id" {
  type = string
}

variable "public_key" {
  type = string
}

variable "esdomain" {
  type    = string
  default = "stresstest"
}

variable "create_es_iam_service_linked_role" {
  type    = bool
  default = false
}

variable "create_esdomain" {
  type    = bool
  default = false
}

variable "datanode_type" {
  type    = string
  default = "m6g.large.elasticsearch"
}

variable "datanode_count" {
  type    = number
  default = 4
}

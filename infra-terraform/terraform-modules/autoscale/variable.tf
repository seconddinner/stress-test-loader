# ###
# # Required variables.
# ###
variable "vpc_id" {
  type = string
}

variable "aws_subnets" {
  type = list(any)

}
variable "subnet_ids" {
  type = list(any)
}


variable "iam_name" {
  type = string   
}

# ###
# # Optional variables. Generally okay to leave at the default.
# ###

# Additional tags to apply to all tagged resources.
variable "extra_tags" {
  type = map(string)
}

variable "stress_test_loader_port" {
  default = 9005
}

variable "cidr_block" {
  type = string
}


variable "min_size" {
  default = 1
}

variable "max_size" {
  default = 3
}

variable "desired_capacity" {
  default = 1
}

variable "instance_type" {
  default = "t2.micro"
}

variable "environment" {
  default = "stress_test_loader"
}

variable "key_name" {
  default = "sd_stresstest"
}

variable "domain_rand" {
  default = ""
}


variable "domain_root" {
  default = "seconddinnertech.net"
}

variable "dns_name" {
  default = "us-west-1"
}

variable "egress_allowed_cidr" {
  default = "0.0.0.0/0"
}

variable "stress_test_loader_allowed_cidr" {

}

variable "ssh_allowed_cidr" {
  type    = list(string)
  default = ["47.154.122.36/32"]
}

variable "prometheus_node_allowed_cidr" {
  default = [ "47.154.122.36/32"]
}

variable "PNS_version" {
  type = string
}

variable "down_scaling_adjustment" {
  default = -1
}

variable "up_scaling_adjustment" {
  default = -1
}

variable "aws_ami_id" {
  type = string
}

variable "user_data" {
}

variable "domain_weight" {
  default = 25
}

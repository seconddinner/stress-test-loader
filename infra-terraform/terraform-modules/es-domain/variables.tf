variable "owner_id" {
  type = string
  default = ""
}

variable "vpc_id" {
  type = string
}

variable "esdomain" {
  type = string
  default = "stresstest"
}

variable "masterusername" {
  type = string
  default = "stresstestadmin"
}

variable "masterpassword" {
  type = string
}

variable "elasticsearch_version" {
  type = string
  default = "7.10"
}

variable "datanode_type" {
  type = string
  default = "m6g.large.elasticsearch"
}

variable "datanode_count" {
  type  = number
  default = 4
}

variable "datanode_size" {
  type  = number
  default = 500
}

variable "masternode_type" {
  type = string
  default = "m6g.large.elasticsearch"
}

variable "masternode_count" {
  type  = number
  default = 3
}

variable "create_iam_service_linked_role" {
  type = bool
  default = true
}


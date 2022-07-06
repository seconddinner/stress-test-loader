variable "cidr_block" {
  description = "default cidr_block"
}

variable "az_count" {
  description = "Number of AZs to cover in a given AWS region"

}

variable "public_ip_on_launch" {
  type    = bool
  default = true
}

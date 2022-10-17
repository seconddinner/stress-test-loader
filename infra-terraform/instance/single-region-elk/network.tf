module "network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
}

data "aws_vpc" "self" {
  id = module.network.aws_vpc_id
}

data "aws_subnets" "self" {
  filter {
    name = "vpc-id"
    values = [module.network.aws_vpc_id]
  }
}
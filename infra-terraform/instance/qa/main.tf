terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.0"

  backend "s3" {

    bucket = "terraform-nvstress-test-sd"
    key    = "qa-stl"
    region = "us-west-2"
  }
}

provider "aws" {
  region = var.aws_region
}

resource "random_id" "id" {
  byte_length = 8
}


locals {
  PNS_version = var.PNS_version != "" ? (var.PNS_version) : (random_id.id.hex)
}

module "network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
}

data "aws_ami" "stl" {
  most_recent = true

  filter {
    name   = "name"
    values = ["stress-test-loader-${var.ami_name}*"]
  }
  owners = ["${var.owner_id}"]
}


locals {
  user_data = templatefile(join("/", tolist([path.module, "user_data.sh"])), {
    stress_test_loader_allowed_cidr = var.stress_test_loader_allowed_cidr
    stress_test_loader_port         = var.stress_test_loader_port
    environment                     = var.environment
  })
}

module "iam" {
  source      = "../../terraform-modules/iam"
  environment = var.environment
}

module "autoscale" {
  source                          = "../../terraform-modules/autoscale"
  vpc_id                          = module.network.aws_vpc_id
  iam_name                        = module.iam.iam_name
  cidr_block                      = var.ntw_cidr_block
  subnet_ids                      = module.network.aws_subnet_list
  aws_subnets                     = module.network.aws_subnets
  min_size                        = var.asg_min
  max_size                        = var.asg_max
  desired_capacity                = var.asg_desired
  key_name                        = var.key_name
  instance_type                   = var.instance_type
  environment                     = var.environment
  stress_test_loader_port         = var.stress_test_loader_port
  PNS_version                     = local.PNS_version
  domain_rand                     = local.PNS_version
  aws_ami_id                      = data.aws_ami.stl.id
  user_data                       = local.user_data
  stress_test_loader_allowed_cidr = var.stress_test_loader_allowed_cidr
  ssh_allowed_cidr                = var.stress_test_loader_allowed_cidr
  extra_tags = {
    "type" = "stress_test_loader"
  }
}


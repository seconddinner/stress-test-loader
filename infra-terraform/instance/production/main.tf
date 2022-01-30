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
    bucket = "terraform-stl-v1"
    key    = "production-stl"
    region = "us-east-1"
  }
}


locals {
  PNS_version = var.PNS_version
}

provider "aws" {
  region  = "us-east-1"
  alias   = "us-east-1"
  profile = "default"
}


provider "aws" {
  alias  = "eu-central-1"
  region = "eu-central-1"
}

provider "aws" {
  alias  = "ap-south-1"
  region = "ap-south-1"
}

provider "aws" {
  alias  = "ap-northeast-2"
  region = "ap-northeast-2"
}

provider "aws" {
  alias  = "ap-southeast-1"
  region = "ap-southeast-1"
}

provider "aws" {
  alias  = "ap-northeast-3"
  region = "ap-northeast-3"
}

provider "aws" {
  alias  = "ap-southeast-2"
  region = "ap-southeast-2"
}

provider "aws" {
  alias  = "ap-northeast-1"
  region = "ap-northeast-1"
}


module "eu-central-1-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.eu-central-1
  }
}

module "ap-south-1-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.ap-south-1
  }
}

module "ap-northeast-2-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.ap-northeast-2
  }
}

module "ap-southeast-1-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.ap-southeast-1
  }
}

module "us-east-1-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.us-east-1
  }
}

module "ami-us-east-1" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.us-east-1
  }
}

module "autoscale-us-east-1" {
  source                  = "../../terraform-modules/autoscale"
  vpc_id                  = module.us-east-1-network.aws_vpc_id
  cidr_block              = var.ntw_cidr_block
  subnet_ids              = module.us-east-1-network.aws_subnet_list
  aws_subnets             = module.us-east-1-network.aws_subnets
  min_size                = var.asg_min
  max_size                = var.asg_max
  desired_capacity        = var.asg_desired
  instance_type           = var.instance_type
  environment             = var.environment
  stress_test_loader_port              = var.stress_test_loader_port
  PNS_version             = local.PNS_version
  down_scaling_adjustment = -(var.asg_min / 2)
  up_scaling_adjustment   = var.asg_min / 2
  domain_rand             = ""
  user_data               = local.user_data
  stress_test_loader_allowed_cidr      = var.stress_test_loader_allowed_cidr
  aws_ami_id              = module.ami-us-east-1.ami_id
  dns_name                = "us-east-1"
  providers = {
    aws = aws.us-east-1
  }
  extra_tags = {
    "type" = "stress_test_loader"
  }
}

module "ami-ap-southeast-1" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.ap-southeast-1
  }
}

module "autoscale-ap-southeast-1" {
  source                  = "../../terraform-modules/autoscale"
  vpc_id                  = module.ap-southeast-1-network.aws_vpc_id
  cidr_block              = var.ntw_cidr_block
  subnet_ids              = module.ap-southeast-1-network.aws_subnet_list
  aws_subnets             = module.ap-southeast-1-network.aws_subnets
  min_size                = var.big_asg_min
  max_size                = var.big_asg_max
  desired_capacity        = var.big_asg_desired
  instance_type           = var.instance_type
  environment             = var.environment
  stress_test_loader_port              = var.stress_test_loader_port
  PNS_version             = local.PNS_version
  down_scaling_adjustment = -(var.big_asg_min / 2)
  up_scaling_adjustment   = var.big_asg_min / 2
  domain_rand             = ""
  user_data               = local.user_data
  stress_test_loader_allowed_cidr      = var.stress_test_loader_allowed_cidr
  aws_ami_id              = module.ami-ap-southeast-1.ami_id
  dns_name                = "ap-southeast-1"
  domain_weight           = 100
  providers = {
    aws = aws.ap-southeast-1
  }
  extra_tags = {
    "type" = "stress_test_loader"
  }
}

module "ami-eu-central-1" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.eu-central-1
  }
}

module "autoscale-eu-central-1" {
  source                  = "../../terraform-modules/autoscale"
  vpc_id                  = module.eu-central-1-network.aws_vpc_id
  cidr_block              = var.ntw_cidr_block
  subnet_ids              = module.eu-central-1-network.aws_subnet_list
  aws_subnets             = module.eu-central-1-network.aws_subnets
  min_size                = var.asg_min
  max_size                = var.asg_max
  desired_capacity        = var.asg_desired
  instance_type           = var.instance_type
  environment             = var.environment
  stress_test_loader_port              = var.stress_test_loader_port
  PNS_version             = local.PNS_version
  down_scaling_adjustment = -(var.asg_min / 2)
  up_scaling_adjustment   = var.asg_min / 2
  domain_rand             = ""
  user_data               = local.user_data
  stress_test_loader_allowed_cidr      = var.stress_test_loader_allowed_cidr
  aws_ami_id              = module.ami-eu-central-1.ami_id
  dns_name                = "eu-central-1"
  providers = {
    aws = aws.eu-central-1
  }
  extra_tags = {
    "type" = "stress_test_loader"
  }
}

module "ami-ap-south-1" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.ap-south-1
  }
}

locals {
  user_data = templatefile(join("/", tolist([path.module, "user_data.sh"])), {
    stress_test_loader_allowed_cidr = var.stress_test_loader_allowed_cidr
    stress_test_loader_port         = var.stress_test_loader_port
    environment        = var.environment
  })
}
# module "autoscale-ap-south-1" {
#   source           = "../../terraform-modules/autoscale"
#   vpc_id           = module.ap-south-1-network.aws_vpc_id
#   cidr_block       = var.ntw_cidr_block
#   subnet_ids       = module.ap-south-1-network.aws_subnet_list
#   aws_subnets      = module.ap-south-1-network.aws_subnets
#   min_size         = 2
#   max_size         = 5
#   desired_capacity = 5
#   instance_type    = var.instance_type
#   environment      = var.environment
#   stress_test_loader_port       = var.stress_test_loader_port
#   PNS_version     = local.PNS_version
#   down_scaling_adjustment = -2
#   up_scaling_adjustment = 2
#   domain_rand = ""
#   aws_ami_id = module.ami-ap-south-1.ami_id
#   dns_name = "ap-south-1"
#   providers = {
#     aws = aws.ap-south-1
#   }
#   extra_tags = {
#     "type" = "stress_test_loader"
#   }
# }

module "ami-ap-northeast-2" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.ap-northeast-2
  }
}

module "autoscale-ap-northeast-2" {
  source                  = "../../terraform-modules/autoscale"
  vpc_id                  = module.ap-northeast-2-network.aws_vpc_id
  cidr_block              = var.ntw_cidr_block
  subnet_ids              = module.ap-northeast-2-network.aws_subnet_list
  aws_subnets             = module.ap-northeast-2-network.aws_subnets
  min_size                = var.asg_min
  max_size                = var.asg_max
  desired_capacity        = var.asg_desired
  instance_type           = var.instance_type
  environment             = var.environment
  stress_test_loader_port              = var.stress_test_loader_port
  PNS_version             = local.PNS_version
  down_scaling_adjustment = -(var.asg_min / 2)
  up_scaling_adjustment   = var.asg_min / 2
  domain_rand             = ""
  user_data               = local.user_data
  stress_test_loader_allowed_cidr      = var.stress_test_loader_allowed_cidr
  aws_ami_id              = module.ami-ap-northeast-2.ami_id
  dns_name                = "ap-northeast-2"
  providers = {
    aws = aws.ap-northeast-2
  }
  extra_tags = {
    "type" = "stress_test_loader"
  }
}

module "ap-northeast-1-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.ap-northeast-1
  }
}

module "ami-ap-northeast-1" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.ap-northeast-1
  }
}

module "autoscale-ap-northeast-1" {
  source                  = "../../terraform-modules/autoscale"
  vpc_id                  = module.ap-northeast-1-network.aws_vpc_id
  cidr_block              = var.ntw_cidr_block
  subnet_ids              = module.ap-northeast-1-network.aws_subnet_list
  aws_subnets             = module.ap-northeast-1-network.aws_subnets
  min_size                = var.asg_min
  max_size                = var.asg_max
  desired_capacity        = var.asg_desired
  instance_type           = var.instance_type
  environment             = var.environment
  stress_test_loader_port              = var.stress_test_loader_port
  PNS_version             = local.PNS_version
  down_scaling_adjustment = -(var.asg_min / 2)
  up_scaling_adjustment   = var.asg_min / 2
  domain_rand             = ""
  user_data               = local.user_data
  stress_test_loader_allowed_cidr      = var.stress_test_loader_allowed_cidr
  aws_ami_id              = module.ami-ap-northeast-1.ami_id
  dns_name                = "ap-northeast-1"
  providers = {
    aws = aws.ap-northeast-1
  }
  extra_tags = {
    "type" = "stress_test_loader"
  }
}


module "ap-southeast-2-network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
  providers = {
    aws = aws.ap-southeast-2
  }
}

module "ami-ap-southeast-2" {
  source        = "../../terraform-modules/ami"
  source_ami_id = var.source_ami_id
  providers = {
    aws = aws.ap-southeast-2
  }
}

module "autoscale-ap-southeast-2" {
  source                  = "../../terraform-modules/autoscale"
  vpc_id                  = module.ap-southeast-2-network.aws_vpc_id
  cidr_block              = var.ntw_cidr_block
  subnet_ids              = module.ap-southeast-2-network.aws_subnet_list
  aws_subnets             = module.ap-southeast-2-network.aws_subnets
  min_size                = var.asg_min
  max_size                = var.asg_max
  desired_capacity        = var.asg_desired
  instance_type           = var.instance_type
  environment             = var.environment
  stress_test_loader_port              = var.stress_test_loader_port
  PNS_version             = local.PNS_version
  down_scaling_adjustment = -(var.asg_min / 2)
  up_scaling_adjustment   = var.asg_min / 2
  domain_rand             = ""
  user_data               = local.user_data
  stress_test_loader_allowed_cidr      = var.stress_test_loader_allowed_cidr
  aws_ami_id              = module.ami-ap-southeast-2.ami_id
  dns_name                = "ap-southeast-2"
  providers = {
    aws = aws.ap-southeast-2
  }
  extra_tags = {
    "type" = "stress_test_loader"
  }
}

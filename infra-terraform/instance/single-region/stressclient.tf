resource "random_id" "id" {
  byte_length = 8
}

locals {
  PNS_version = var.PNS_version != "" ? (var.PNS_version) : (random_id.id.hex)
}

data "aws_ami" "stl" {
  most_recent = true

  filter {
    name   = "name"
    values = ["stress-test-loader-${var.ami_name}*"]
  }
  filter {
    name   = "architecture"
    values = ["${var.clientarch}"]
  }
  owners = ["${var.owner_id}"]
}

locals {
  user_data = templatefile(join("/", tolist([path.module, "user_data.sh"])), {
    stress_test_loader_allowed_cidr = var.stress_test_loader_allowed_cidr
    stress_test_loader_port         = var.stress_test_loader_port
    environment                     = var.environment
    masterusername                  = var.masterusername
    masterpassword                  = var.masterpassword
    es_endpoint                     = var.create_esdomain ? module.es-domain[0].es_endpoint : data.aws_elasticsearch_domain.es[0].endpoint
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
  public_key                      = var.public_key
  extra_tags = {
    "type" = "stress_test_loader"
  }
}


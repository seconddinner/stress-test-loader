data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
    account_id = data.aws_caller_identity.current.account_id
    owner_id = length(var.owner_id) > 0 ? var.owner_id : local.account_id
}

module "network" {
  source     = "../../terraform-modules/vpc"
  cidr_block = var.ntw_cidr_block
  az_count   = var.az_count
}

data "aws_vpc" "self" {
  id = module.network.aws_vpc_id
}

data "aws_subnet_ids" "self" {
  vpc_id = module.network.aws_vpc_id
}


resource "aws_security_group" "es_sg" {
  description = "controls direct access to opensearch instance"
  vpc_id      = data.aws_vpc.self.id
  name        = "snap-stresstest-es-sg"

  ingress {
    protocol  = "tcp"
    from_port = 443
    to_port   = 443

    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

}

resource "aws_iam_service_linked_role" "es" {
  aws_service_name = "es.amazonaws.com"
}

resource "aws_elasticsearch_domain" "es" {
  domain_name           = var.esdomain
  elasticsearch_version = "7.10"
  domain_endpoint_options {
    enforce_https = true
    tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
  }
  cluster_config {
    instance_count = 4
    instance_type          = "m6g.large.elasticsearch"
    zone_awareness_enabled = false
    dedicated_master_enabled = true
    dedicated_master_count = 3
    dedicated_master_type = "m6g.large.elasticsearch"

  }
  ebs_options {
    ebs_enabled = true
    volume_size = 500
    volume_type = "gp3"
  }

  #vpc_options {
  #  subnet_ids = [
  #    tolist(data.aws_subnet_ids.self.ids)[0],
  #    tolist(data.aws_subnet_ids.self.ids)[1],
  #  ]
#
  #  security_group_ids = [aws_security_group.es_sg.id]
  #}

  advanced_options = {
    "rest.action.multi.allow_explicit_index" = "true"
  }

  node_to_node_encryption {
    enabled = true
  }
  encrypt_at_rest {
    enabled = true
  }
  advanced_security_options {
    enabled = true
    internal_user_database_enabled = true
    master_user_options {
        master_user_name = "stresstestadmin"
        master_user_password = var.masterpassword

    }

  }
  auto_tune_options {
    desired_state = "ENABLED"
    rollback_on_disable = "NO_ROLLBACK"

  }

  access_policies = <<CONFIG
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "es:*",
            "Principal": "*",
            "Effect": "Allow",
            "Resource": "arn:aws:es:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:domain/${var.esdomain}/*"
        }
    ]
}
CONFIG

  tags = {
    Domain = "stresstestdomain"
  }

  depends_on = [aws_iam_service_linked_role.es]
}

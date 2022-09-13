data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
    account_id = data.aws_caller_identity.current.account_id
    owner_id = length(var.owner_id) > 0 ? var.owner_id : local.account_id
}

data "aws_vpc" "self" {
  id = var.vpc_id
}

data "aws_subnets" "self" {
  filter {
    name = "vpc-id"
    values = [data.aws_vpc.self.id]
  }
}

locals {
  subnet_ids = data.aws_subnets.self.ids
}

resource "aws_security_group" "es_sg" {
  description = "controls direct access to opensearch instance"
  vpc_id      = data.aws_vpc.self.id
  name        = "snap-es-${var.esdomain}-sg"

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
  count            = var.create_iam_service_linked_role ? 1 : 0
  aws_service_name = "es.amazonaws.com"
}

resource "aws_elasticsearch_domain" "es" {
  domain_name           = var.esdomain
  elasticsearch_version = var.elasticsearch_version
  domain_endpoint_options {
    enforce_https = true
    tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
  }
  cluster_config {
    instance_count = var.datanode_count
    instance_type          = var.datanode_type
    zone_awareness_enabled = false
    dedicated_master_enabled = true
    dedicated_master_count = var.masternode_count
    dedicated_master_type = var.masternode_type

  }
  ebs_options {
    ebs_enabled = true
    volume_size = var.datanode_size
    volume_type = "gp3"
  }

  # use the vpc_options when es domain needs to accepts the logs inside the vpc
  #vpc_options {
  #  subnet_ids = [
  #    tolist(local.subnet_ids)[0],
  #    tolist(local.subnet_ids)[1],
  #  ]
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
        master_user_name = var.masterusername
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

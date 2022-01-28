resource "aws_alb" "main" {
  name    = "${format("%s-%s", var.aws_alb_name, var.deployname)}"
  subnets = ["${var.aws_subnet_id}"]

  security_groups = ["${aws_security_group.lb_sg.id}"]

  tags {
    terraform_module = "s3_alb"
  }
}

resource "aws_alb_listener" "front_end" {
  load_balancer_arn = "${aws_alb.main.id}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = "${var.aws_alb_listener_certificate_arn}"

  default_action {
    target_group_arn = "${var.aws_alb_target_group_ecs_target_id}"
    type             = "forward"
  }
}

resource "aws_ecs_service" "ecs_service_main" {
  name            = "${var.ecs_service_main_name}"
  cluster         = "${var.aws_ecs_cluster_main_id}"
  task_definition = "${var.aws_ecs_task_definition_ecs_task_def_arn}"
  desired_count   = "${var.aws_ecs_service_desired_count}"
  iam_role        = "${aws_iam_role.ecs_service.name}"

  load_balancer {
    target_group_arn = "${var.aws_alb_target_group_ecs_target_id}"
    container_name   = "${var.s3_alb_container_name}"
    container_port   = "${var.container_port}"
  }

  depends_on = [
    "aws_iam_role_policy.ecs_service",
    "aws_alb_listener.front_end",
  ]
}

resource "aws_iam_role" "ecs_service" {
  name = "${format("%s_%s", var.ecs_iam_service_role, var.deployname)}"

  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ecs_service" {
  name = "${var.ecs_service_iam_role_policy}"
  role = "${aws_iam_role.ecs_service.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:Describe*",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:DeregisterTargets",
        "elasticloadbalancing:Describe*",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:RegisterTargets"
      ],
      "Resource": "*"
    },{
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject"
            ],
            "Resource": "arn:aws:s3:::${var.artifactdomainname}/*"
   },
   {
    "Effect": "Allow",
    "Action": "kms:Decrypt",
    "Resource": "${var.kms_decrypt_key_arn}"
    }
  ]
}
EOF
}

resource "aws_s3_bucket" "b" {
  bucket        = "${var.domainname}"
  acl           = "public-read"
  force_destroy = true

  website {
    index_document = "index.html"

    routing_rules = <<EOF
[{
    "Condition": {
        "HttpErrorCodeReturnedEquals": "404"
    },
    "Redirect": {
       "HostName": "${var.domainname}",
        "ReplaceKeyPrefixWith": "#!/"
    }
}]
EOF
  }

  tags {
    terraform_module = "s3_alb"
  }
}

resource "aws_s3_bucket_policy" "b" {
  bucket = "${aws_s3_bucket.b.id}"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::${aws_s3_bucket.b.bucket}/*"
            ]
        }
    ]
}
POLICY
}

resource "aws_security_group" "lb_sg" {
  description = "controls access to the application ELB"

  vpc_id = "${var.aws_vpc_id}"
  name   = "${var.aws_security_group_lb_sg_name}"

  ingress {
    protocol  = "tcp"
    from_port = 443
    to_port   = 443

    cidr_blocks = [
      "${var.admin_cidr_ingress}",
      "0.0.0.0/0",
    ]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"

    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }

  tags {
    terraform_module = "s3_alb"
  }
}

resource "aws_route53_record" "api" {
  zone_id = "${var.route53_zone_id}"
  name    = "${format("api.%s", var.domainname)}"

  type = "A"

  alias {
    name                   = "${aws_alb.main.dns_name}"
    zone_id                = "${aws_alb.main.zone_id}"
    evaluate_target_health = true
  }
}

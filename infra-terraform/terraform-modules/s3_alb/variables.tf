variable "ecs_service_main_name" {
  description = ""
  default     = "ecs_service_main_name_default"
}

variable "aws_ecs_service_desired_count" {
  description = ""
}

variable "ecs_iam_service_role" {
  description = ""
  default     = "ecs_iam_service_role_default"
}

variable "aws_ecs_cluster_main_id" {
  description = ""
}

variable "ecs_service_iam_role_policy" {
  description = ""
  default     = "ecs_service_iam_role_policy_default"
}

variable "aws_iam_instance_profile_name" {
  description = ""
  default     = "aws_iam_instance_profile_name_default"
}

variable "aws_alb_listener_certificate_arn" {
  description = "alb listener certificate arn"
}

variable "s3_alb_container_name" {
  description = ""
}

variable "container_port" {
  description = ""
}

variable "aws_alb_target_group_ecs_target_id" {
  description = ""
}

variable "aws_ecs_task_definition_ecs_task_def_arn" {
  description = ""
}

variable "domainname" {
  description = ""
}

variable "route53_zone_id" {
  description = ""
}

variable "aws_subnet_id" {
  type        = "list"
  description = ""
}

variable "aws_alb_name" {
  default = "aws-alb-default"
}

variable "aws_security_group_lb_sg_name" {
  default = "aws_security_group_lb_sg_default"
}

variable "admin_cidr_ingress" {
  description = ""
}

variable "aws_vpc_id" {
  description = ""
}

variable "kms_decrypt_key_arn" {
  description = ""
}

variable "deployname" {
  description = ""
}

variable "artifactdomainname" {}

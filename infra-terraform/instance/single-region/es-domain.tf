module "es-domain" {
  source = "../../terraform-modules/es-domain"
  esdomain = var.esdomain
  vpc_id = module.network.aws_vpc_id
  masterpassword = var.masterpassword
  create_iam_service_linked_role = var.create_es_iam_service_linked_role
}
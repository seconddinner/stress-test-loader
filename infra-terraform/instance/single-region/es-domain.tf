# create new es domain when 'create_esdomain' is set to 'true' explicitly
module "es-domain" {
  source = "../../terraform-modules/es-domain"
  count = var.create_esdomain ? 1 : 0
  esdomain = var.esdomain
  vpc_id = module.network.aws_vpc_id
  masterusername = var.masterusername
  masterpassword = var.masterpassword
  create_iam_service_linked_role = var.create_es_iam_service_linked_role
}

# leverage the existing es domain by default
data "aws_elasticsearch_domain" "es" {
  count = var.create_esdomain ? 0 : 1
  domain_name = var.esdomain
}
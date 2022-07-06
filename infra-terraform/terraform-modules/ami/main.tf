resource "aws_ami_copy" "stl" {
  name              = var.name
  source_ami_id = var.source_ami_id
  source_ami_region = var.source_ami_region
}
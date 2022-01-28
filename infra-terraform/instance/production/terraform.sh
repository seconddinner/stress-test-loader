#!/usr/bin/env bash
set -x
cd cicd/terraform/aws/instance/production
terraform init
echo $image_name
echo $ami_id
terraform apply  -auto-approve -var="PNS_version=$image_name" -var="source_ami_id=$ami_id" || true
terraform apply  -auto-approve -var="PNS_version=$image_name" -var="source_ami_id=$ami_id" || true
terraform apply  -auto-approve -var="PNS_version=$image_name" -var="source_ami_id=$ami_id"
cd ../../../../..
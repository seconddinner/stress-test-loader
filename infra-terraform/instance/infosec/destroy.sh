#!/usr/bin/env bash
cd cicd/terraform/aws/instance/infosec
terraform init
terraform destroy  -auto-approve -var="PNS_version=$image_name" 
cd ../../../../..


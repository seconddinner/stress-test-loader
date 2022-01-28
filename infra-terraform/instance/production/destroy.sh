#!/usr/bin/env bash
cd cicd/terraform/aws/instance/production
terraform init
terraform destroy  -auto-approve -var="PNS_version=$image_name" 
cd ../../../../..


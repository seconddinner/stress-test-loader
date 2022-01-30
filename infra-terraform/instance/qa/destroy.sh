#!/usr/bin/env bash
cd cicd/terraform/aws/instance/qa-us-west-1
sed -i "s/___NEED_SED___/qa-$PNS_version/g" qa.tfbackend
terraform init -reconfigure -backend-config=qa.tfbackend
terraform destroy  -auto-approve -var="PNS_version=$PNS_version" || true
terraform destroy  -auto-approve -var="PNS_version=$PNS_version" || true
terraform destroy  -auto-approve -var="PNS_version=$PNS_version"
cd ../../../../..
aws s3api delete-object --bucket terraform-hmg-v3  --key "qa-$PNS_version"

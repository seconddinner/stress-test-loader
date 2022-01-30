#!/usr/bin/env bash
set -x
cd cicd/terraform/aws/instance/qa-us-west-1
sed -i "s/___NEED_SED___/qa-$PNS_version/g" qa.tfbackend
terraform init -reconfigure -backend-config=qa.tfbackend 
if [[ $latestAMIVersion == "true" ]] 
then 
    terraform apply  -auto-approve -var="PNS_version=$PNS_version" 
else 
    terraform apply  -auto-approve -var="PNS_version=$PNS_version" -var="ami_name=$PNS_version"
fi
export ami_id_q=$(terraform output ami_id)
export ami_id=`echo $ami_id_q |tr -d '"'`
export stress_test_loader_qa_url=$(terraform output stress_test_loader_url)
export stress_test_loader_qa_url=`echo $stress_test_loader_qa_url | tr -d '"'`
echo $ami_id
echo $stress_test_loader_qa_url
cd ../../../../..
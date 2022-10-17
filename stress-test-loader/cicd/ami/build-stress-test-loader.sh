#!/usr/bin/env bash
set -x 
export STL_version=`git describe --tags --long`
echo "${STL_version}"
bash build.sh
# packer init cicd/packer
# packer fmt cicd/packer
# packer build -var="PNS_version=stress-test-loader-${STL_version}" -machine-readable cicd/packer/ubuntu-stress-test-loader.pkr.hcl | tee build.log
# export PROXY_NODE_SELFCHECK_AMI_ID=`grep 'artifact,0,id' build.log | cut -d, -f6 | cut -d: -f2`
set +x
# [Second Dinner](https://seconddinner.com/work-together-at-second-dinner/) stress-test-loader

This is the [Second Dinner](https://seconddinner.com/work-together-at-second-dinner/) stress-test-loader. It is a small golang application for executing stress tests, plus some infra things (packer, terraform, etc) to deploy it to the cloud. We have open-sourced this code to contribute to development community. 

Currently, this setup targets AWS, but it can be ported other clouds if needed. 

Directory structure:

* stress-test-loader (golang service that can load any stress-test-client, plus packer templates for creating AMIs)

* infra-terraform (terraform config for deploying AMIs to create stress-test-loader ec2 instances)

## Requirements

1. [golang](https://go.dev/doc/install).
1. [hashicorp packer](https://www.packer.io/downloads).
1. [configure your AWS environment keys](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).
1. [AWS Cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
1. [hashicorp terraform](https://www.terraform.io/downloads).

## stress-test-loader

To take advantage of [AWS arm64 offering](https://aws.amazon.com/ec2/graviton/), we are building the golang executable and packer AMI image all in arm64 format.

The following example demonstrates building the stress-test-loader ami.

### Generate AMI

1. ``` cd stress-test-loader ```
1. ``` source cicd/ami/build-stress-test-loader.sh ```
1. If everything worked according to plan, you will see message like ``` AMIs were created: ami-XXXXXXXXXXXXXXX ```

## Infra-terraform

Once you have created an AWS AMI for stress-test, you can use Infra-terraform to create EC2 instance and create as many EC2 instances as your AWS account allows.

### terraform stress infrastructure 

1. make a copy of ```infra-terraform/instance/single-region``` in ```infra-terraform/instance/```, we use ```infra-terraform/instance/single-region-test``` as an example.
1. replace `backend "s3"` in ```infra-terraform/instance/single-region-test/main.tf``` with your own backend.
1. need following variables: [public_key](https://www.techrepublic.com/article/how-to-view-your-ssh-keys-in-linux-macos-and-windows/), your aws account id as [owner_id](https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html), your machines public ip as [stress_test_loader_allowed_cidr](https://ifconfig.me/)
1. update variables `infra-terraform/instance/single-region-test/variable.tf`.
1. cd ```infra-terraform/instance/single-region-test```
1. ```terraform init```
1. ```terraform apply```
1. If everything worked according to plan, you will see message like ``` Apply complete! Resources: XX added, XX changed, XX destroyed ```

## Running the stress test using stress-test-loader

### Create stresstest client json
* build your stress test client as an arm64 executable, this can be a directory with libraries and one entry executable. The executable can take any number of environment variable as configuration. We are going to use ```stress-test-client``` as an example
* ```tar cvzf stress-test-client```
* copy the tgz file to an S3 bucket ```aws s3 cp  stress-test-client.tgz   s3://stress-test-client-s3/stress-test-client.tgz```
* build a stress-test-loader config json. For example `stresstest.json` 
```
{
    "s3": "stress-test-client-s3",
    "s3key": "stress-test-client.tgz",
    "loadtestExec": "StressTest.Cli",
    "envVariableList": [
        {
            "EnvName": "TargetEnvironment",
            "EnvValue": "EnvValue"
        },
        {
            "EnvName": "SomeOtherEnv",
            "EnvValue": "us-west-2"
        }
    ]
}
```
* export all of the ec2 instances' public IP addresses ```aws ec2 describe-instances --region us-west-2   --query 'Reservations[*].Instances[*].{"public_ip":PublicIpAddress}' --filters Name=instance-state-name,Values=running --output json     > /tmp/IP.json```
* run stress test ``` cd stress-test-loader/client; go run main.go stresstest.json /tmp/IP.json```
* if you gave an ssh public key, you can ssh into the ec2 instance and check its systemd service log ```journalctl -f -u stress*```


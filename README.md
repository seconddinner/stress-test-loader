# [Second Dinner](https://seconddinner.com/) stress-test-loader

This is [Second Dinner](https://seconddinner.com/) stress-test-loader. We want to open source this code to contribute to development community. 

Current setup is targeted AWS, but it can be ported other clouds if needed. 

Current dir structure have:

* stress-test-loader (golang service that can load any stress-test-client, the dir also create AWS ami that allow user to create hundreds of loadtest ec2 instances)

* infra-terraform. terraform code that allow folks to AWS ami to create stress-test-loader ec2 instances)

## requirement:

1. [golang](https://go.dev/doc/install).
1. [hashicorp packer](https://www.packer.io/downloads).
1. [configure your AWS environment keys](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).
1. [AWS Cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
1. [hashicorp terraform](https://www.terraform.io/downloads).

## stress-test-loader

To take advantage of [AWS arm64 offering](https://aws.amazon.com/ec2/graviton/), we are building the golang executable and packer ami image all in arm64 format.

Here is brief process of building stress-test-loader ami.

### generate ami

1. ``` cd stress-test-loader ```
1. ``` source cicd/ami/build-stress-test-loader.sh ```
1. If everything worked according to plan, you will see message like ``` AMIs were created: ami-XXXXXXXXXXXXXXX ```

## Infra-terraform

Once you created an AWS ami for stress-test, you can use Infra-terraform to create EC2 instance and create as many EC2 instances as your AWS accounts allowed.

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
* build your stress test client to an arm64 executable, this can be a directory with libraries and one entry executable. Executable can take any number environment variable as configuration. We are going to use ```stress-test-client``` as an example
* ```tar cvzf stress-test-client```
* copy the tgz file to a S3 bucket ```aws s3 cp  stress-test-client.tgz   s3://stress-test-client-s3/stress-test-client.tgz```
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
* export all of the ec2 instances public IP address ```aws ec2 describe-instances --region us-west-2   --query 'Reservations[*].Instances[*].{"public_ip":PublicIpAddress}' --filters Name=instance-state-name,Values=running --output json     > /tmp/IP.json```
* run stress test ``` cd stress-test-loader/client; go run main.go stresstest.json /tmp/IP.json```
* if you gave a ssh public key, you can ssh into the ec2 instance and check log ```journal-ctl -f -u stress*```


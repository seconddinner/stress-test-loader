# [Second Dinner](https://seconddinner.com/) stress-test loader

This is how Second dinner stress test our platform. We want to open source this code to contribute to development community. 

Current dir structure have:

* stress-test-loader (golang service that can load any stress-test-client, the dir also create AWS ami that allow user to create hundreds of loadtest ec2 instances)

* infra-terraform. terraform code that allow folks to AWS ami to create stress-test-loader ec2 instances)

## requirement:

1. [golang](https://go.dev/doc/install).
1. [hashicorp packer](https://www.packer.io/downloads).
1. [configure your AWS environment keys](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).
1. [golang](https://go.dev/doc/install).
1. [AWS Cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
1. [hashicorp terraform](https://www.terraform.io/downloads).

## stress-test-loader

To take advantage of AWS arm64 offering, we are building the golang executable and packer ami image all in arm64 format.

Here is brief process of building stress-test-loader ami.

### generate ami

1. ``` cd stress-test-loader ```
1. ``` source cicd/ami/build-stress-test-loader.sh ```
1. If everything worked according to plan, you will see message like ``` AMIs were created: ami-XXXXXXXXXXXXXXX ```

## Infra-terraform

Once you created an AWS ami for stress-test, you can use Infra-terraform to create EC2 instance and create as many EC2 instances as your AWS accounts allowed.

### terraform stress infrastructure 

1. make a copy of ```infra-terraform/instance/qa``` in ```infra-terraform/instance/```, we use ```infra-terraform/instance/qa2``` as an example.
1. update variable.tf according to your aws account and local public IP settings.
1. cd ```infra-terraform/instance/qa2```
1. ```terraform init```
1. ```terraform apply```
1. If everything worked according to plan, you will see message like ``` Apply complete! Resources: XX added, XX changed, XX destroyed ```

## Running the stress test using stresstest loader

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


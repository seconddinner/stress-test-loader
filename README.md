# [Second Dinner](https://seconddinner.com/work-together-at-second-dinner/) stress-test-loader

This is the [Second Dinner](https://seconddinner.com/work-together-at-second-dinner/) stress-test-loader. It is a small golang application for executing stress tests (and arguably any executable) on the cloud, plus some infra things (packer, pulumi, etc) for deployment. We have open-sourced this code to contribute to development community. 

Currently, this setup targets AWS, but it can be ported other clouds if needed. 

Directory structure:

* stress-test-loader (golang service that can load any stress-test-client, plus packer templates for creating AMIs)

* infra-pulumi (.NET Pulumi code for deploying AMIs to create stress-test-loader ec2 instances)

* build-container (Docker container for running the stress test loader in GitHub Actions)

* simple-stress-test-client (A simple stress test client that pings any host)

## Requirements

1. [golang](https://go.dev/doc/install).
1. [hashicorp packer](https://www.packer.io/downloads).
1. [configure your AWS environment keys](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).
1. [AWS Cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
1. [Pulumi](https://www.pulumi.com/docs/install/).
1. [OpenSSL](https://www.openssl.org/source/).

## stress-test-loader

To take advantage of [AWS arm64 offering](https://aws.amazon.com/ec2/graviton/), we are building the golang executable and packer AMI image all in arm64 format.

The following example demonstrates building the stress-test-loader ami.

### Generate SSL certificates

To secure the GRPC connection between the client and server of the stress test loader, we use SSL to authenticate clients. Please follow the steps below to generate the certificates.
1. ``` cd stress-test-loader ```
1. Generate a private key and self-signed certificate for a CA ``` ./gen_ca.sh ```
1. Generate certificates for the server and the client ``` ./gen_cert.sh ```

### Generate AMI

1. ``` cd stress-test-loader ```
1. ``` source cicd/ami/build-stress-test-loader.sh ```
1. If everything worked according to plan, you will see message like ``` AMIs were created: ami-XXXXXXXXXXXXXXX ```

## infra-pulumi

Once you have created an AWS AMI for stress-test, you can use infra-pulumi to create EC2 instance and create as many EC2 instances as your AWS account allows.

### pulumi stress infrastructure 

1. Need following variables: 
* [public_key](https://www.techrepublic.com/article/how-to-view-your-ssh-keys-in-linux-macos-and-windows/): your ssh public key; 
* [stress_test_loader_allowed_cidr](https://ifconfig.me/): your machine's public ip; 
* [s3_client_bucket_name](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingBucket.html): the name of your AWS S3 bucket to store the stress test client executable; 
* [s3_log_bucket_name](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingBucket.html): the name of your AWS S3 bucket to store the logs; 
* [desired_capacity](https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-capacity-limits.html): the number of ec2 instances to create in a region; 
* [regions](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html): the AWS regions to create ec2 instances, separated by commas (e.g., "us-east-1, us-west-2"); 
1. Update variables `infra-pulumi/Infra.Pulumi/Infra.Pulumi/Config.cs`.
1. ```cd infra-pulumi/Infra.Pulumi/Infra.Pulumi```
1. Set up the infrastructure ```dotnet run --project-name stress-test-loader-pulumi --stack-name dev```
1. If everything worked according to plan, you will see message like 
``` 
Diagnostics:
  pulumi:pulumi:Stack (stress-test-loader-pulumi-dev):
    Downloading provider: aws

Resources:
    + 51 created

Duration: 4m20s 
```
1. the public IP addresses of all the ec2 instances will be stored in `/tmp/IP.json`
1. to destroy the infrastructure after the stress test: ```dotnet run --project-name stress-test-loader-pulumi --stack-name dev --destroy```

## Running the stress test using stress-test-loader

### Build your stress test client and upload it to S3
1. build your stress test client as an arm64 executable, this can be a directory with libraries and one entry executable. The executable can take any number of environment variable as configuration. We are going to use ```simple-stress-test-client``` as an example
1. ```cd simple-stress-test-client/StressTest```
1. ```dotnet publish -r linux-arm64 --self-contained true -c Release```
1. ```cd bin/Release```
1. ```cd "$(ls -d */ | head -n 1)"```
1. ```cd linux-arm64/publish```
1. ```tar czf /tmp/simple-stress-test-client.tgz ./```
1. copy the tgz file to an S3 bucket ```aws s3 cp  /tmp/simple-stress-test-client.tgz   s3://stress-test-client-s3/simple-stress-test-client.tgz```

### Create a stress test client configuration json
1. Build a stress-test-loader config json. For example `stresstest.json` 
```
{
    "s3": "stress-test-client-s3",
    "s3key": "simple-stress-test-client.tgz",
    "loadtestExec": "StressTest",
    "envVariableList": [
        {
            "EnvName": "num_pings",
            "EnvValue": "10"
        },
        {
            "EnvName": "ping_interval",
            "EnvValue": "500"
        },
        {
            "EnvName": "host",
            "EnvValue": "https://www.google.com/"
        }
    ]
}
```

### Run the stress test
1. ```cd stress-test-loader/client```
1. Run stress test ```go run main.go stresstest.json /tmp/IP.json```
1. If you gave an ssh public key, you can ssh into the ec2 instance and check its systemd service log ```journalctl -f -u stress*```
1. If you are running our simple-stress-test-client, you can check the log ```cat /tmp/stresstest-log```
1. To check the status of the stress tests (running or finished) ```go run main.go -p /tmp/IP.json```
1. To stop the stress tests ```go run main.go -s /tmp/IP.json```


# [Second Dinner](https://seconddinner.com/work-together-at-second-dinner/) stress-test-loader

This is the [Second Dinner](https://seconddinner.com/work-together-at-second-dinner/) stress-test-loader. It is a small golang application for executing stress tests (and arguably any executable) on the cloud, plus some infra things (packer, pulumi, etc) for deployment. We have open-sourced this code to contribute to development community. 

Currently, this setup targets AWS, but it can be ported other clouds if needed. 

[Note] We are transitioning from terraform to pulumi for ease of maintain

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
* [stress_test_loader_allowed_cidr](https://ifconfig.me/): your machine's public IP; when running in GitHub Actions, this should be the GitHub runner's public IP (feel free to check out our workflow);
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
1. You can let the stress test loader client to wait for the stress tests to finish by specifying a maximum time to wait
1. If you gave an ssh public key, you can ssh into the ec2 instance and check its systemd service log ```journalctl -f -u stress*```
1. If you are running our simple-stress-test-client, you can check the log ```cat /tmp/stresstest-log```
1. To check the status of the stress tests (running or finished) ```go run main.go -p /tmp/IP.json```
1. To stop the stress tests ```go run main.go -s /tmp/IP.json```


## How to use our GitHub Actions workflows?
We provided three reusable GitHub Actions workflows, namely [_stress-test-packer-build.yaml](https://github.com/seconddinner/stress-test-loader/blob/main/.github/workflows/_stress-test-packer-build.yaml), [_pulumi-set-up.yaml](https://github.com/seconddinner/stress-test-loader/blob/main/.github/workflows/_pulumi-set-up.yaml), and [_run-stress-test.yaml](https://github.com/seconddinner/stress-test-loader/blob/main/.github/workflows/_run-stress-test.yaml). You can fork this repo and call these workflows from your repo.

### [_stress-test-packer-build.yaml](https://github.com/seconddinner/stress-test-loader/blob/main/.github/workflows/_stress-test-packer-build.yaml)
This workflow builds the stress test loader, generates certificates for the SSL connection between stress test loader clients and servers, and creates an AMI using Packer.
To call this workflow, you will need to provide the following secrets:
1. [TARGET_ACCOUNT](https://docs.aws.amazon.com/IAM/latest/UserGuide/FindingYourAWSId.html): Your AWS Account ID.
1. [CA_KEY](https://gist.github.com/Soarez/9688998): The private key of your certificate authority (CA) used in SSL. You may generate a new key by `openssl genrsa -out example.org.key 2048`

### [_pulumi-set-up.yaml](https://github.com/seconddinner/stress-test-loader/blob/main/.github/workflows/_pulumi-set-up.yaml)
This workflow sets up the infrastructure on AWS for stress test clients. The infrastructure includes S3 buckets (may already exist), IAM roles and policies, VPC, and AutoScaling Groups in one or more regions.
To call this workflow, you will need to provide the following input:
1. [REGIONS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html): The AWS regions to create ec2 instances, separated by commas (e.g., "us-east-1, us-west-2").
1. [DESIRED_CAPACITY](https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-capacity-limits.html): The number of ec2 instances to create in a region.

And the following secrets:
1. [GITHUB_ACTION_PULUMI_ACCESS_TOKEN](https://www.pulumi.com/docs/pulumi-cloud/access-management/access-tokens/): Your Pulumi Cloud access token.
1. [STRESSTESTLOADER_S3_CLIENT_BUCKET_NAME](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html): The name of your S3 bucket that stores the stress test client.
1. [STRESSTESTLOADER_S3_LOG_BUCKET_NAME](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html): The name of your S3 bucket that will store the stress test logs.
1. [TARGET_ACCOUNT](https://docs.aws.amazon.com/IAM/latest/UserGuide/FindingYourAWSId.html): Your AWS Account ID.
1. [CA_KEY](https://gist.github.com/Soarez/9688998): The private key of your certificate authority (CA) used in SSL. This should be the same one you used for the Packer build.

This workflow will produce one output:
1. [EC2_IP](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-instance-addressing.html): The public IP addresses of the EC2 instances created by this workflow.

### [_run-stress-test.yaml](https://github.com/seconddinner/stress-test-loader/blob/main/.github/workflows/_run-stress-test.yaml)
This workflow starts the stress tests on the launched EC2 instances. It will wait until all stress tests are finished or the specified timeout. Optionally, it can ssh into the EC2 instances and fetch the logs. In the end, it will destroy the infrastructure.
To call this workflow, you will need to provide the following input:
1. [GET_RESULTS](https://github.com/seconddinner/stress-test-loader/blob/0ecba75183a16c55b9b34c3389fc09b886ae3243/.github/workflows/_run-stress-test.yaml#LL104C8-L104C8): A boolean. If set, the stress test loader will ssh into the EC2 instances and print the last 100 lines of `/tmp/stress-test-log` to the GitHub Actions console.
1. [EC2_IP](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-instance-addressing.html): The public IP addresses of the EC2 instances output by `_pulumi-set-up.yaml`.
1. [REGIONS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html): The AWS regions to create ec2 instances, separated by commas (e.g., "us-east-1, us-west-2").
1. [DESIRED_CAPACITY](https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-capacity-limits.html): The number of ec2 instances to create in a region.

And the following secrets:
1. [GITHUB_ACTION_PULUMI_ACCESS_TOKEN](https://www.pulumi.com/docs/pulumi-cloud/access-management/access-tokens/): Your Pulumi Cloud access token.
1. [STRESSTESTLOADER_S3_CLIENT_BUCKET_NAME](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html): The name of your S3 bucket that stores the stress test client.
1. [STRESSTESTLOADER_S3_LOG_BUCKET_NAME](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html): The name of your S3 bucket that will store the stress test logs.
1. [TARGET_ACCOUNT](https://docs.aws.amazon.com/IAM/latest/UserGuide/FindingYourAWSId.html): Your AWS Account ID.
1. [STRESS_TEST_JSON](https://jsonlint.com/): The config json for the stress test loader. Please see the section `Create a stress test client configuration json` in this README.
1. [STRESS_TEST_TOTAL_TIME](https://www.loadview-testing.com/blog/how-long-should-a-website-load-test-run/#:~:text=As%20a%20rule%20of%20thumb,test%20conducted%20and%20your%20goals.): The maximum time (in seconds) to run the stress test. After the specified time, any running stress test will be stopped by the loader.
1. [SSH_PRIVATE_KEY](https://www.techrepublic.com/article/how-to-view-your-ssh-keys-in-linux-macos-and-windows/): Your ssh public key.
1. [CA_KEY](https://gist.github.com/Soarez/9688998): The private key of your certificate authority (CA) used in SSL. This should be the same one you used for the Packer build.

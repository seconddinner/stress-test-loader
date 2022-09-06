terraform {
    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 4"
        }
    }
    required_version = "~> 1.0"
    backend "s3" {
        bucket = "terraform-nvstress-test-sd"
        key = "qa-stl-opensearch"
        region = "us-west-2"
        profile = "aws_aws-cube-nvstressclient-bd-us_app.luoweichao"

    }
}

provider "aws" {
    region = "us-west-2"
    profile = "aws_aws-cube-nvstressclient-bd-us_app.luoweichao"
}
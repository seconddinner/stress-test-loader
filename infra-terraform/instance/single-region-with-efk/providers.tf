terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.0"

  backend "s3" {
    bucket = "terraform-nvstress-test-sd"
    key    = "luoweichao-stl"
    region = "us-west-2"
    profile = "aws_aws-cube-nvstressclient-bd-us_app.luoweichao"
  }
}

provider "aws" {
  region = var.aws_region
  profile = "aws_aws-cube-nvstressclient-bd-us_app.luoweichao"
}
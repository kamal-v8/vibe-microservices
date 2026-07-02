terraform {
  # required_version = "1.5.0"

  backend "s3" {
    bucket       = "amzn-bucket-for-pulse"
    key          = "pulse/platform/terraform.tfstate"
    region       = "us-east-1"
    use_lockfile = true
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}



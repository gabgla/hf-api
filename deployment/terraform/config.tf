terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.62.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "4.1.0"
    }
  }

  backend "s3" {
    bucket = locals.statefile_bucket
    key    = locals.statefile_key
    region = var.aws_region
  }
}

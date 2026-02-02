// AWS Config
variable "aws_region" {
  description = "The AWS region to deploy resources in"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "The AWS CLI profile to use"
  type        = string
  default     = "default"
}

// Deployment Config
variable "namespace" {
  description = "The namespace to prefix resource names with"
  type        = string
  default     = "hfapi"
}

variable "environment" {
  description = "The deployment environment (e.g., dev, staging, prod). 'prod' and 'live' will not be prefixed."
  type        = string
}

variable "tags" {
  description = "A map of tags to apply to resources"
  type        = map(string)
  default     = {}
}

// VPC Config
variable "create_vpc" {
  description = "Whether to create the VPC"
  type        = bool
  default     = true
}

variable "create_public_subnet" {
  description = "Whether to create the public subnet"
  type        = bool
  default     = true
}

variable "vpc_cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
  description = "CIDR block for the public subnet. Leave empty to calculate"
  type        = string
  default     = null
}

// DNS Config
variable "hosted_zone_name" {
  description = "Name of the Hosted Zone, hosted in AWS"
  type        = string
}

variable "api_subdomain" {
  description = "Subdomain name to be combined with the hosted zone, namespace and environment"
  type        = string
  default     = "api"
}

// Lambda Config
variable "lambda_memory_size" {
  description = "The amount of memory (in MB) to allocate to the Lambda function"
  type        = number
  default     = 128
}

variable "lambda_timeout" {
  description = "The timeout (in seconds) for the Lambda function"
  type        = number
  default     = 2
}

variable "lambda_environment_variables" {
  description = "A map of environment variables to set for the Lambda function"
  type        = map(string)
  default     = {}
}

variable "lambda_reserved_concurrent_executions" {
  description = "The number of concurrent executions reserved for the Lambda function"
  type        = number
  default     = 10
}

// API Gateway Config

variable "api_gateway_throttling_burst_limit" {
  description = "The API Gateway throttling burst limit"
  type        = number
  default     = 100
}

variable "api_gateway_throttling_rate_limit" {
  description = "The API Gateway throttling rate limit"
  type        = number
  default     = 50
}

variable "cors_allow_origins" {
  description = "CORS allowed origins"
  type        = list(string)
  default     = ["*"]
}

variable "cors_allow_methods" {
  description = "CORS allowed methods"
  type        = list(string)
  default     = ["GET", "OPTIONS"]
}

variable "cors_allow_headers" {
  description = "CORS allowed headers"
  type        = list(string)
  default     = ["Content-Type", "Authorization"]
}

variable "cors_expose_headers" {
  description = "CORS headers to expose"
  type        = list(string)
  default     = []
}

variable "cors_max_age" {
  description = "CORS preflight cache duration in seconds"
  type        = number
  default     = 3600
}

// Logging and Monitoring Config
variable "enable_cloudwatch_logging" {
  description = "Enable CloudWatch logging for the API Gateway"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "The number of days to retain CloudWatch logs"
  type        = number
  default     = 7
}

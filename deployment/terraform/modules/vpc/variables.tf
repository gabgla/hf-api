variable "prefix" {
  description = "Prefix for all VPC resources"
  type        = string
}

variable "create_vpc" {
  description = "Whether to create the VPC"
  type        = bool
}

variable "create_public_subnet" {
  description = "Whether to create the public subnet"
  type        = bool
}

variable "vpc_cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "public_subnet_cidr" {
  description = "CIDR block for the public subnet"
  type        = string
}

variable "tags" {
  description = "A map of tags to apply to resources"
  type        = map(string)
}

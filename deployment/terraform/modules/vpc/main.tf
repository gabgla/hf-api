// VPC
resource "aws_vpc" "new_vpc" {
  count = var.create_vpc ? 1 : 0

  cidr_block           = var.vpc_cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = var.tags
}

data "aws_vpc" "existing_vpc" {
  count = var.create_vpc ? 0 : 1

  cidr_block = var.vpc_cidr_block
}

locals {
  vpc_id = var.create_vpc ? aws_vpc.new_vpc[0].id : data.aws_vpc.existing_vpc[0].id
}

// Subnet
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "public" {
  count = var.create_public_subnet ? 1 : 0

  vpc_id                  = local.vpc_id
  cidr_block              = var.public_subnet_cidr != null ? var.public_subnet_cidr : cidrsubnet(var.vpc_cidr_block, 8, 1)
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.prefix}public-subnet"
  }
}

data "aws_subnet" "existing_public_subnet" {
  count = var.create_public_subnet ? 0 : 1

  vpc_id     = local.vpc_id
  cidr_block = var.public_subnet_cidr
}

locals {
  public_subnet_id = var.create_public_subnet ? aws_subnet.public[0].id : data.aws_subnet.existing_public_subnet[0].id
}

// Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = local.vpc_id

  tags = {
    Name = "${var.prefix}-igw"
  }
}

resource "aws_route_table" "public" {
  vpc_id = local.vpc_id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.prefix}-public-rt"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = local.public_subnet_id
  route_table_id = aws_route_table.public.id
}


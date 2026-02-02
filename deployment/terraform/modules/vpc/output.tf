output "vpc_id" {
  description = "The ID of the VPC"
  value       = local.vpc_id
}

output "public_subnet_id" {
  description = "The ID of the public subnet"
  value       = local.public_subnet_id
}

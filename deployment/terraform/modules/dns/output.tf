output "zone_id" {
  description = "The ID of the hosted zone"
  value       = data.aws_route53_zone.root.zone_id
}

output "certificate_arn" {
  description = "The ARN of the validated ACM certificate"
  value       = aws_acm_certificate_validation.api.certificate_arn
}

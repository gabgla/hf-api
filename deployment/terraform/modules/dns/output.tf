output "zone_id" {
  description = "The ID of the hosted zone"
  value       = aws_route53_zone.main.zone_id
}

output "certificate_arn" {
  description = "The ARN of the ACM certificate"
  value       = aws_acm_certificate.main.arn
}

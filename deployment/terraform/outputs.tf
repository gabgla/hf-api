output "artifact_s3_bucket" {
  description = "The S3 bucket for artifacts"
  value       = aws_s3_bucket.artifacts.id
}

output "artifact_s3_key" {
  description = "The S3 key for the Lambda artifact"
  value       = aws_s3_object.lambda_placeholder.key
}

output "api_lambda_function_name" {
  description = "The name of the API Lambda function"
  value       = aws_lambda_function.api.function_name
}

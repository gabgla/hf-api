output "artifact_s3_artifact_key" {
  description = "The S3 artifact key"
  value       = aws_s3_object.lambda_placeholder.key
}

output "api_lambda_function_name" {
  description = "The name of the API Lambda function"
  value       = aws_lambda_function.api.function_name
}

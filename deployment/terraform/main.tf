# -----------------------------------------------------------------------------
# Dependencies
# -----------------------------------------------------------------------------

# module "vpc" {
#   source = "./modules/vpc"

#   prefix               = local.resource_prefix
#   vpc_cidr_block       = var.vpc_cidr_block
#   public_subnet_cidr   = var.public_subnet_cidr
#   create_vpc           = var.create_vpc
#   create_public_subnet = var.create_public_subnet
#   tags                 = var.tags
# }

module "dns" {
  source = "./modules/dns"

  hosted_zone_name = var.hosted_zone_name
  api_subdomain    = var.api_subdomain
  tags             = var.tags
}

# -----------------------------------------------------------------------------
# Lambda Function
# -----------------------------------------------------------------------------

resource "aws_lambda_function" "api" {
  function_name = local.function_name
  description   = "HFAPI Lambda Function"
  runtime       = "provided.al2023"
  handler       = "bootstrap"

  filename         = "${workspace.root_module.directory}/../artifacts/bootstrap.zip"
  source_code_hash = filebase64sha256("${workspace.root_module.directory}/../artifacts/bootstrap.zip")

  role        = aws_iam_role.lambda_exec.arn
  memory_size = var.lambda_memory_size
  timeout     = var.lambda_timeout

  environment {
    variables = var.lambda_environment_variables
  }

  tags = var.tags
}

# -----------------------------------------------------------------------------
# IAM Role for Lambda
# -----------------------------------------------------------------------------

resource "aws_iam_role" "lambda_exec" {
  name = "${local.function_name}-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# resource "aws_iam_role_policy_attachment" "lambda_vpc" {
#   count      = var.vpc_config != null ? 1 : 0
#   role       = aws_iam_role.lambda_exec.name
#   policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
# }

# -----------------------------------------------------------------------------
# CloudWatch Log Group
# -----------------------------------------------------------------------------

resource "aws_cloudwatch_log_group" "api" {
  count = var.enable_cloudwatch_logging ? 1 : 0
  name              = "/aws/lambda/${local.function_name}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# -----------------------------------------------------------------------------
# API Gateway (HTTP API)
# -----------------------------------------------------------------------------

resource "aws_apigatewayv2_api" "api" {
  name          = local.function_name
  protocol_type = "HTTP"
  description   = "HFAPI HTTP API Gateway"

  cors_configuration {
    allow_origins     = "*"
    allow_methods     = "OPTIONS,GET,POST,PUT,DELETE,PATCH"
    allow_headers     = "*"
    expose_headers    = "*"
    max_age           = "3600"
    allow_credentials = false
  }

  tags = var.tags
}

resource "aws_apigatewayv2_stage" "api" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format = jsonencode({
      requestId        = "$context.requestId"
      ip               = "$context.identity.sourceIp"
      requestTime      = "$context.requestTime"
      httpMethod       = "$context.httpMethod"
      routeKey         = "$context.routeKey"
      status           = "$context.status"
      responseLength   = "$context.responseLength"
      integrationError = "$context.integrationErrorMessage"
    })
  }

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "api_gateway" {
  count = var.enable_cloudwatch_logging ? 1 : 0
  name              = "/aws/api-gateway/${local.function_name}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.api.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "catch_all" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

# -----------------------------------------------------------------------------
# Custom Domain
# -----------------------------------------------------------------------------

resource "aws_apigatewayv2_domain_name" "api" {
  domain_name = var.api_subdomain

  domain_name_configuration {
    certificate_arn = module.dns.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "api" {
  api_id      = aws_apigatewayv2_api.http.id
  domain_name = aws_apigatewayv2_domain_name.api.id

  stage       = aws_apigatewayv2_stage.api.id
}

resource "aws_route53_record" "api_a" {
  zone_id = module.dns.zone_id
  name    = var.api_subdomain
  type    = "A"

  alias {
    name                   = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_aaaa" {
  zone_id = data.aws_route53_zone.root.zone_id
  name    = var.api_subdomain
  type    = "AAAA"

  alias {
    name                   = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

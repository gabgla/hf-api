#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="${SCRIPT_DIR}/../.."
TF_DIR="${SCRIPT_DIR}/../terraform"
BUILD_DIR="${ROOT_DIR}/build/lambda"

cd "${ROOT_DIR}"

# Build the Lambda artifact
make build-for-lambda

# Create deployment package
cd "${BUILD_DIR}"
zip -r lambda.zip bootstrap index.bleve

# Get terraform outputs
cd "${TF_DIR}"
S3_BUCKET=$(terraform output -raw artifact_s3_bucket)
S3_KEY=$(terraform output -raw artifact_s3_key)
FUNCTION_NAME=$(terraform output -raw api_lambda_function_name)

# Upload to S3
aws s3 cp "${BUILD_DIR}/lambda.zip" "s3://${S3_BUCKET}/${S3_KEY}"

# Deploy to Lambda
aws lambda update-function-code \
    --function-name "${FUNCTION_NAME}" \
    --s3-bucket "${S3_BUCKET}" \
    --s3-key "${S3_KEY}"

echo "Deployed ${FUNCTION_NAME} successfully"

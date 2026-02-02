#!/bin/bash

set -euo pipefail

# Terraform statefile bootstrap script
# Usage: ./bootstrap.sh <namespace> <environment>

if [[ $# -lt 2 ]]; then
    echo "Usage: $0 <namespace> <environment>"
    exit 1
fi

NAMESPACE="$1"
ENVIRONMENT="$2"
STATE_BUCKET="tf-state-${NAMESPACE}-${ENVIRONMENT}"
REGION="${AWS_REGION:-us-east-1}"

echo "Bootstrapping Terraform state for namespace=$NAMESPACE, environment=$ENVIRONMENT"

# Create S3 bucket for state
aws s3api create-bucket \
    --bucket "$STATE_BUCKET" \
    --region "$REGION" \
    $([ "$REGION" != "us-east-1" ] && echo "--create-bucket-configuration LocationConstraint=$REGION") \
    2>/dev/null || true

# Enable versioning
aws s3api put-bucket-versioning \
    --bucket "$STATE_BUCKET" \
    --versioning-configuration Status=Enabled

# Enable encryption
aws s3api put-bucket-encryption \
    --bucket "$STATE_BUCKET" \
    --server-side-encryption-configuration '{
        "Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]
    }'

# Block public access
aws s3api put-public-access-block \
    --bucket "$STATE_BUCKET" \
    --public-access-block-configuration \
    "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

echo "✓ Terraform state bootstrap complete"
echo "State bucket: $STATE_BUCKET"

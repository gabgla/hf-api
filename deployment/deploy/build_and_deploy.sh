BUILD_DIR=build/lambda

# Setup
make setup

# Build the Lambda artifact
make build-for-lambda

# Create deployment package
zip ${BUILD_DIR}/lambda.zip -j ${BUILD_DIR}/bootstrap ${BUILD_DIR}/index.bleve

# Upload to S3
aws s3 cp ${BUILD_DIR}/lambda.zip s3://$(terraform output -raw artifact_s3_artifact_key)

# Deploy to Lambda
aws lambda update-function-code \
    --function-name $(terraform output -raw api_lambda_function_name) \
    --s3-bucket $(terraform output -raw artifact_s3_artifact_key) \
    --s3-key lambdas/hfapi/lambda.zip

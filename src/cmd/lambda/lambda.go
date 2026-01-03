package main

import (
	_ "embed"
	"hf-api/src/internal/app/api"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	handler := api.NewRouterHandler()
	adapter := httpadapter.New(handler)
	lambda.Start(adapter.ProxyWithContext)
}

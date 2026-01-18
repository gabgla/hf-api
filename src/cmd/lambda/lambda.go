package main

import (
	"log"

	"hf-api/src/internal/app/api"
	"hf-api/src/internal/data"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	err := data.LoadDB()
	if err != nil {
		log.Fatal(err)
	}

	handler := api.NewRouterHandler()
	adapter := httpadapter.New(handler)
	lambda.Start(adapter.ProxyWithContext)
}

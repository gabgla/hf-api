package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"hfapi-app/src/internal/app/api"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestAdapter struct {
	*events.APIGatewayProxyRequest
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s (requestId=%s)", req.HTTPMethod, req.Path, req.RequestContext.RequestID)

	res := api.Router(ctx, &RequestAdapter{APIGatewayProxyRequest: &req})
	return lambdaResponse(res)
}

func lambdaResponse(resp *api.APIResponse) (events.APIGatewayProxyResponse, error) {
	lambdaResponse := events.APIGatewayProxyResponse{
		StatusCode: resp.Code,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
		},
	}

	var content any

	if resp.Error != nil {
		content = resp.Error
	} else if resp.Content != nil {
		content = resp.Content
	}

	if resp.Content != nil {
		buf, err := json.Marshal(content)

		if err != nil {
			log.Printf("Error marshaling JSON response: %v", err)
			lambdaResponse.StatusCode = http.StatusInternalServerError
			return lambdaResponse, err
		}

		lambdaResponse.Body = string(buf)
		lambdaResponse.Headers["Content-Type"] = "application/json"
	}

	return lambdaResponse, nil
}

func main() {
	lambda.Start(router)
}

func (r *RequestAdapter) Method() string {
	return r.APIGatewayProxyRequest.HTTPMethod
}

func (r *RequestAdapter) URL() *url.URL {
	url, _ := url.Parse(r.APIGatewayProxyRequest.Path)
	return url
}

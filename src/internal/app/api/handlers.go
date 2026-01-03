package api

import (
	"context"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func health(ctx context.Context, req *http.Request) *APIResponse {
	return &APIResponse{
		Code: http.StatusOK,
		Content: &HealthResponse{
			Status: "ok",
		},
	}
}

func search(ctx context.Context, req *http.Request) *APIResponse {
	return &APIResponse{
		Code: http.StatusOK,
		Content: &HealthResponse{
			Status: "ok",
		},
	}
}

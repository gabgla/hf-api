package api

import (
	"context"
	"net/http"
	"net/url"
)

type APIHandler func(ctx context.Context, req APIRequest) *APIResponse

type APIResponse struct {
	Code    int
	Error   error
	Content any
}

type RouteKey struct {
	Method string
	Path   string
}

type HealthResponse struct {
	Status string `json:"status"`
}

type APIRequest interface {
	Method() string
	URL() *url.URL
}

var routes = map[RouteKey]APIHandler{
	{"GET", "/health"}: health,
	{"GET", "/search"}: search,
}

func Router(ctx context.Context, req APIRequest) *APIResponse {
	if r, ok := routes[RouteKey{req.Method(), req.URL().Path}]; ok {
		return r(ctx, req)
	}

	for k := range routes {
		if k.Path == req.URL().Path {
			return &APIResponse{
				Code:  http.StatusMethodNotAllowed,
				Error: wrapError("Method not allowed", nil),
			}
		}
	}

	return &APIResponse{
		Code:  http.StatusNotFound,
		Error: wrapError("Not found", nil),
	}
}

func Routes() map[RouteKey]APIHandler {
	return routes
}

// Handlers

func health(ctx context.Context, req APIRequest) *APIResponse {
	return &APIResponse{
		Code: http.StatusOK,
		Content: &HealthResponse{
			Status: "ok",
		},
	}
}

func search(ctx context.Context, req APIRequest) *APIResponse {
	return &APIResponse{
		Code: http.StatusOK,
		Content: &HealthResponse{
			Status: "ok",
		},
	}
}

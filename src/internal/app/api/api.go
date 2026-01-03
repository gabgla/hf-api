package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type APIHandler func(ctx context.Context, req *http.Request) *APIResponse

type APIResponse struct {
	Code    int
	Error   error
	Content any
}

type HealthResponse struct {
	Status string `json:"status"`
}

type APIRequest interface {
	Method() string
	URL() *url.URL
}

var routes = map[string]APIHandler{
	"GET /health": health,
	"GET /search": search,
}

func NewRouterHandler() http.Handler {
	mux := http.NewServeMux()

	for pattern, handler := range routes {
		mux.Handle(pattern, envelopeMiddleware(handler))
	}

	return mux
}

func envelopeMiddleware(handler APIHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		res := handler(ctx, r)

		w.WriteHeader(res.Code)

		var body any

		if res.Error != nil {
			body = res.Error
		} else if res.Content != nil {
			body = res.Content
		}

		if body != nil {
			w.Header().Set("Content-Type", "application/json")
			buf, _ := json.Marshal(body)
			_, _ = w.Write(buf)
		}
	})
}

// Handlers

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

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

type APIRequest interface {
	Method() string
	URL() *url.URL
}

var routes = map[string]APIHandler{
	"GET /health":       health,
	"GET /cards/search": search,
}

func NewRouterHandler() http.Handler {
	mux := http.NewServeMux()

	for pattern, handler := range routes {
		mux.Handle(pattern, envelopeMiddleware(handler))
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/v1/", http.StripPrefix("/v1", mux))

	return rootMux
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

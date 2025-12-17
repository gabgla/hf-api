package main

import (
	_ "embed"
	"encoding/json"
	"hfapi-app/src/internal/app/api"
	"log"
	"net/http"
	"net/url"
)

type RequestAdapter struct {
	*http.Request
}

func main() {
	runServer()
}

func runServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := api.Router(r.Context(), &RequestAdapter{Request: r})

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

	addr := ":8080"
	log.Printf("Running local server on %s ...", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func (r *RequestAdapter) Method() string {
	return r.Request.Method
}

func (r *RequestAdapter) URL() *url.URL {
	return r.Request.URL
}

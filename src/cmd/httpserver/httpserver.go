package main

import (
	"encoding/json"
	"hfapi-app/src/internal/app/api"
	"log"
	"net/http"
	"net/url"
)

type RequestWrapper struct {
	*http.Request
}

func main() {
	runServer()
}

func runServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := api.Router(r.Context(), &RequestWrapper{Request: r})

		w.WriteHeader(res.Code)

		if res.Error != nil {
			buf, _ := json.Marshal(res.Error)
			_, _ = w.Write(buf)
			return
		}

		if res.Content != nil {
			buf, _ := json.Marshal(res.Error)
			_, _ = w.Write(buf)
		}
	})

	addr := ":8080"
	log.Printf("Running local server on %s ...", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func (r *RequestWrapper) Method() string {
	return r.Request.Method
}

func (r *RequestWrapper) URL() *url.URL {
	return r.Request.URL
}

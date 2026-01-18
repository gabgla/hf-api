package main

import (
	_ "embed"
	"hf-api/src/internal/app/api"
	"hf-api/src/internal/data"
	"log"
	"net/http"
)

func main() {
	err := data.LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	runServer()
}

func runServer() {
	handler := api.NewRouterHandler()
	log.Default().Println("Starting HTTP server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

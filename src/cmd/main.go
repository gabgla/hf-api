package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"hfapi-app/src/pkg/cards"
	"hfapi-app/src/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ErrorResponse is a generic error payload.
type ErrorResponse struct {
	Message string `json:"message"`
}

// HealthResponse is returned by the /health endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

const max_sides = 4 // Current maximum number of sides for any card

//go:embed database.json
var dbJSON []byte

var db []cards.Card
var dbIndex map[string]*cards.Card

func loadDB() {
	var dbParsed cards.Root
	if err := json.Unmarshal(dbJSON, &dbParsed); err != nil {
		log.Fatalf("Failed to unmarshal database JSON: %v", err)
	}

	db = normaliseDB(&dbParsed)

	// Build index: Name -> Set
	dbIndex = make(map[string]*cards.Card)
	for _, c := range db {
		if c.Name == "" {
			continue
		}

		dbIndex[strings.ToLower(c.Name)] = &c
	}
}

func normaliseDB(db *cards.Root) []cards.Card {
	result := []cards.Card{}
	for _, c := range db.Data {
		card := cards.Card{
			Name:          c.Name,
			Creator:       c.Creator,
			Set:           c.Set,
			Legality:      c.ConstructedLegality,
			Rulings:       c.Rulings,
			ManaValue:     utils.Coalesce(c.CMC, 0),
			Colors:        strings.Split(c.Colors, ";"),
			Sides:         parseSides(c),
			Tags:          strings.Split(c.Tags, ";"),
			Tokens:        c.Tokens,
			ComponentOf:   c.ComponentOf,
			IsActualToken: c.IsActualToken,
			SmallAltImage: c.SmallAltImage,
		}
		result = append(result, card)
	}

	return result
}

func parseSides(c cards.CardEntry) []cards.Side {
	sides := []cards.Side{}

	if c.Cost == nil {
		c.Cost = []*string{nil, nil, nil, nil}
	}

	if c.Supertypes == nil {
		c.Supertypes = []*string{nil, nil, nil, nil}
	}

	if c.CardTypes == nil {
		c.CardTypes = []*string{nil, nil, nil, nil}
	}

	if c.Subtypes == nil {
		c.Subtypes = []*string{nil, nil, nil, nil}
	}

	if c.Power == nil {
		c.Power = []*any{nil, nil, nil, nil}
	}

	if c.Toughness == nil {
		c.Toughness = []*any{nil, nil, nil, nil}
	}

	if c.Loyalty == nil {
		c.Loyalty = []*any{nil, nil, nil, nil}
	}

	if c.TextBox == nil {
		c.TextBox = []*string{nil, nil, nil, nil}
	}

	if c.FlavorText == nil {
		c.FlavorText = []*string{nil, nil, nil, nil}
	}

	for i := range max_sides {
		if c.Cost[i] == nil &&
			c.Supertypes[i] == nil &&
			c.CardTypes[i] == nil &&
			c.Subtypes[i] == nil &&
			c.Power[i] == nil &&
			c.Toughness[i] == nil &&
			c.Loyalty[i] == nil &&
			c.TextBox[i] == nil &&
			c.FlavorText[i] == nil {
			continue
		}

		// For numeric fields, attempt to parse integers
		// This allows for mathematical comparisons later
		powerStr, powerInt := getStringAndNumber(c.Power[i])
		toughnessStr, toughnessInt := getStringAndNumber(c.Toughness[i])
		loyaltyStr, loyaltyInt := getStringAndNumber(c.Loyalty[i])

		if len(c.FlavorText) < i+1 {
			fmt.Println(c.Name, c.FlavorText)
		}

		side := cards.Side{
			Cost:              utils.Coalesce(c.Cost[i], ""),
			Supertypes:        strings.Split(utils.Coalesce(c.Supertypes[i], ""), ";"),
			CardTypes:         strings.Split(utils.Coalesce(c.CardTypes[i], ""), ";"),
			Subtypes:          strings.Split(utils.Coalesce(c.Subtypes[i], ""), ";"),
			Power:             powerInt,
			Toughness:         toughnessInt,
			Loyalty:           loyaltyInt,
			PowerOriginal:     powerStr,
			ToughnessOriginal: toughnessStr,
			LoyaltyOriginal:   loyaltyStr,
			TextBox:           utils.Coalesce(c.TextBox[i], ""),
			FlavorText:        utils.Coalesce(c.FlavorText[i], ""),

			Tags:   strings.Split(c.Tags, ";"),
			Tokens: c.Tokens,
		}

		sides = append(sides, side)
	}

	return sides
}

func getStringAndNumber(field *any) (string, *int) {
	var fieldStr string
	var fieldInt *int

	if p := field; p != nil {
		if v, ok := (*p).(string); ok {
			fieldStr = v
		}
	}

	if n, err := strconv.Atoi(fieldStr); err == nil {
		fieldInt = &n
	}

	return fieldStr, fieldInt
}

// router is the main Lambda handler.
func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s (requestId=%s)", req.HTTPMethod, req.Path, req.RequestContext.RequestID)

	if db == nil {
		loadDB()
	}

	switch req.Path {
	case "/health":
		if req.HTTPMethod != http.MethodGet {
			return clientError(http.StatusMethodNotAllowed)
		}
		return jsonResponse(http.StatusOK, HealthResponse{Status: "ok"})

	case "/items":
		switch req.HTTPMethod {
		case http.MethodGet:
			// Example: list items
			items := []map[string]interface{}{
				{"id": "1", "name": "Foo"},
				{"id": "2", "name": "Bar"},
			}
			return jsonResponse(http.StatusOK, map[string]interface{}{
				"items": items,
			})

		case http.MethodPost:
			// Example: create item from JSON body
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
				log.Printf("Error decoding body: %v", err)
				return clientError(http.StatusBadRequest)
			}

			// Here you'd normally save to DB, etc.
			// For now, just echo back with a fake ID.
			payload["id"] = "123"
			return jsonResponse(http.StatusCreated, payload)

		default:
			return clientError(http.StatusMethodNotAllowed)
		}
	case "/search":
		if req.HTTPMethod != http.MethodGet {
			return clientError(http.StatusMethodNotAllowed)
		}

		name, ok := req.QueryStringParameters["name"]
		if !ok {
			return jsonResponse(http.StatusOK, []cards.Card{})
		}

		card, ok := dbIndex[strings.ToLower(name)]
		if !ok {
			return jsonResponse(http.StatusOK, []cards.Card{})
		}

		return jsonResponse(http.StatusOK, []*cards.Card{card})

	default:
		return clientError(http.StatusNotFound)
	}
}

// jsonResponse marshals the body and returns a JSON API Gateway response.
func jsonResponse(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
	js, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		return serverError()
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",            // basic CORS
			"Access-Control-Allow-Headers": "Content-Type", // tweak as needed
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
		},
		Body: string(js),
	}, nil
}

// clientError wraps a standard HTTP status with a JSON error response.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return jsonResponse(status, ErrorResponse{
		Message: http.StatusText(status),
	})
}

// serverError is a generic 500 error response.
func serverError() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"message":"internal server error"}`,
	}, nil
}

// --- LOCAL SERVER MODE ---

// toAPIGatewayRequest converts an http.Request into a minimal APIGatewayProxyRequest.
func toAPIGatewayRequest(r *http.Request) (events.APIGatewayProxyRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return events.APIGatewayProxyRequest{}, err
	}
	// allow the body to be read again if needed
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	qs := map[string]string{}
	for k, v := range r.URL.Query() {
		qs[k] = strings.Join(v, ",")
	}

	return events.APIGatewayProxyRequest{
		Path:                            r.URL.Path,
		QueryStringParameters:           qs,
		MultiValueQueryStringParameters: r.URL.Query(),
		HTTPMethod:                      r.Method,
		Headers:                         headers,
		Body:                            string(bodyBytes),
	}, nil
}

// runLocalServer starts an HTTP server that uses the Lambda router under the hood.
func runLocalServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req, err := toAPIGatewayRequest(r)
		if err != nil {
			log.Printf("Error converting request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"conversion error"}`))
			return
		}

		resp, err := router(context.Background(), req)
		if err != nil {
			log.Printf("Handler error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"handler error"}`))
			return
		}

		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write([]byte(resp.Body))
	})

	addr := ":8080"
	log.Printf("Running local server on %s ...", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func main() {
	if os.Getenv("LOCAL") == "1" {
		runLocalServer()
		return
	}

	lambda.Start(router)
}

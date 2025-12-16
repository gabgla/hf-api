package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"hfapi-app/src/internal/app/api"
	"hfapi-app/src/pkg/cards"
	"hfapi-app/src/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ErrorResponse is a generic error payload.
type ErrorResponse struct {
	Message string `json:"message"`
}

type RequestWrapper struct {
	*events.APIGatewayProxyRequest
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

func getStringAndNumber(field *any) (*string, *int) {
	var fieldStr string
	var fieldInt *int

	if p := field; p != nil {
		if v, ok := (*p).(string); ok {
			fieldStr = v
		}
	} else {
		return nil, nil
	}

	if n, err := strconv.Atoi(fieldStr); err == nil {
		fieldInt = &n
	}

	return &fieldStr, fieldInt
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s (requestId=%s)", req.HTTPMethod, req.Path, req.RequestContext.RequestID)

	if db == nil {
		loadDB()
	}

	res := api.Router(ctx, &RequestWrapper{APIGatewayProxyRequest: &req})
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

func (r *RequestWrapper) Method() string {
	return r.APIGatewayProxyRequest.HTTPMethod
}

func (r *RequestWrapper) URL() *url.URL {
	url, _ := url.Parse(r.APIGatewayProxyRequest.Path)
	return url
}

package api

import (
	"context"
	"fmt"
	"hf-api/src/internal/data"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// /health

func health(ctx context.Context, req *http.Request) *APIResponse {
	return &APIResponse{
		Code: http.StatusOK,
		Content: &HealthResponse{
			Status: "ok",
		},
	}
}

// /cards/search

const epsilon = 10e-3 // We don't need super high precision for this

type Filter struct {
	Key      string
	Operator string
	Value    string
}

var KnownTokens = map[string][]string{
	"name":        {"name", "n"},
	"colors":      {"colors", "c", "color"},
	"mv":          {"mv", "cmc"},
	"mana":        {"mana", "m"},
	"identity":    {"identity", "id"},
	"type_line":   {"type_line", "type", "t"},
	"oracle":      {"oracle", "o"},
	"flavor_text": {"flavor_text", "ft", "flavor", "flavortext"},

	"power":           {"power", "pow"},
	"toughness":       {"toughness", "tou", "tough"},
	"power_toughness": {"power_toughness", "pt", "powtou"},
	"loyalty":         {"loyalty", "loy"},

	"devotion": {"devotion"},
	"produces": {"produces"},

	"set":     {"set", "s", "edition", "e"},
	"tags":    {"tags", "tag"},
	"creator": {"creator", "author"},

	"format": {"format", "f"},
	"banned": {"banned", "ban"},
}

func search(ctx context.Context, req *http.Request) *APIResponse {
	query := req.URL.Query().Get("q")

	if query == "" {
		return &APIResponse{
			Code:  http.StatusBadRequest,
			Error: &APIError{Message: "Empty search"},
		}
	}

	filters := extractFitersFromQuery(query)

	// Fallback: if no filters were extracted, treat entire query as name search
	if len(filters) == 0 {
		filters = []Filter{
			{Key: "name", Operator: ":", Value: query},
		}
	}

	searchRequest := bleve.NewSearchRequest(buildBleveQuery(filters))
	searchRequest.Fields = []string{"*"}
	searchRequest.Size = 10
	// searchRequest.SortBy([]string{"-_score"})

	results, err := data.Index.SearchInContext(ctx, searchRequest)

	if err != nil {
		return &APIResponse{
			Code:  http.StatusInternalServerError,
			Error: &APIError{Message: "Search failed"},
		}
	}

	return &APIResponse{
		Code:    http.StatusOK,
		Content: results,
	}
}

func buildBleveQuery(filters []Filter) *query.BooleanQuery {
	query := bleve.NewBooleanQuery()

	for _, filter := range filters {
		if filter.Operator == ":" {
			matchQuery := bleve.NewMatchQuery(filter.Value)
			matchQuery.SetField(filter.Key)
			query.AddMust(matchQuery)
			continue
		}

		var num, min, max float64
		var err error

		if num, err = strconv.ParseFloat(filter.Value, 64); err != nil {
			continue
		}

		switch filter.Operator {
		case "=":
			min = num - epsilon
			max = num + epsilon
		case ">":
			min = num + epsilon
			max = math.Inf(1)
		case "<":
			min = math.Inf(-1)
			max = num - epsilon
		case ">=":
			min = num - epsilon
			max = math.Inf(1)
		case "<=":
			min = math.Inf(-1)
			max = num + epsilon
		}

		rangeQuery := bleve.NewNumericRangeQuery(&min, &max)
		rangeQuery.SetField(filter.Key)
		query.AddMust(rangeQuery)
	}

	return query
}

func extractFitersFromQuery(query string) []Filter {
	// Extract tokens from query https://scryfall.com/docs/syntax
	re := regexp.MustCompile(`(\w+)([:=><]=?)(?:"([^"]+)"|([^"\s]+))(?:\s|$)`)
	matches := re.FindAllStringSubmatch(query, -1)

	filters := []Filter{}
	for _, match := range matches {
		if len(match) < 4 {
			fmt.Println("malformed filter:", match[0])
			// malformed filter; skip
			continue
		}

		keyAlias := strings.TrimSpace(match[1])
		key, ok := TokenAliasMap[keyAlias]

		if !ok {
			fmt.Println("unknown filter key:", keyAlias)
			// unknown filter key; skip
			continue
		}

		op := strings.TrimSpace(match[2])
		value := strings.TrimSpace(match[3])

		if value == "" {
			value = strings.TrimSpace(match[4])
		}

		filters = append(filters, Filter{
			Key:      key,
			Operator: op,
			Value:    value,
		})
	}

	return filters
}

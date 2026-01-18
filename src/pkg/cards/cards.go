package cards

import (
	"math"
	"strconv"
	"strings"
)

type Card struct {
	// Identifiers
	Name    string `json:"name"`
	Creator string `json:"creator"`
	Set     string `json:"set"`

	// Legality
	Legality []string `json:"legality"`
	Rulings  string   `json:"rulings"`

	// Characteristics
	ManaValue *float64 `json:"mv"`
	Colors    []string `json:"colors"`
	Sides     []Side   `json:"sides"` // most cards have 1 side; some have up to 4

	// Refs
	Tags   []string `json:"tags"`
	Tokens []Token  `json:"-"`

	// Additional fields
	ComponentOf   *string `json:"component_of"`
	IsActualToken *bool   `json:"is_actual_token"`
	SmallAltImage *string `json:"small_alt_image"`

	ManaValueOriginal *string `json:"mv_original"`
}

type Side struct {
	Cost       string   `json:"cost"`
	Supertypes []string `json:"supertypes"`
	CardTypes  []string `json:"card_types"`
	Subtypes   []string `json:"subtypes"`
	ManaValue  float64  `json:"mv"`
	Power      *float64 `json:"power"`
	Toughness  *float64 `json:"toughness"`
	Loyalty    *float64 `json:"loyalty"`
	TextBox    string   `json:"textbox"`
	FlavorText string   `json:"flavor_text"`

	// Original fields in case numeric parsing doesn't work
	ManaValueOriginal *string `json:"mv_original"`
	PowerOriginal     *string `json:"power_original"`
	ToughnessOriginal *string `json:"toughness_original"`
	LoyaltyOriginal   *string `json:"loyalty_original"`
}

type Token struct {
	Name      string `json:"Name"`
	Power     string `json:"Power"`
	Toughness string `json:"Toughness"`
	Type      string `json:"Type"`
	Image     string `json:"Image"`
}

var ManaTokensWihoutIdentity = map[string]struct{}{
	"H": {},
	"X": {},
	"Y": {},
	"Z": {},
}

var KnownManaTokens = map[string]float64{
	"W": 1,
	"U": 1,
	"B": 1,
	"R": 1,
	"G": 1,
	"P": 1, // Purple, not Phyrexian
	"C": 1,
	"X": 0,
	"Y": 0,
	"Z": 0,
	"H": 0, // Phyrexian. No such thing as "pay life only" without alternative mana cost
}

func ParseManaValueInt(cost string) (int, error) {
	v, err := ParseManaValue(cost)
	return int(v), err
}

func ParseManaValueFloat(cost string) (float64, error) {
	return ParseManaValue(cost)
}

func ParseManaValue(cost string) (float64, error) {
	cost = strings.TrimSpace(cost)

	total := 0.0
	invalidSyntax := false
	open := false
	symbols := []string{}

	currentSymbol := ""

	for _, r := range cost {
		switch r {
		case '{':
			if open {
				invalidSyntax = true
			}
			open = true
		case '/':
			if !open {
				invalidSyntax = true
			}

			if currentSymbol != "" {
				symbols = append(symbols, currentSymbol)
				currentSymbol = ""
			} else {
				invalidSyntax = true
			}
		case '}':
			if !open {
				invalidSyntax = true
			}
			open = false

			if currentSymbol != "" {
				symbols = append(symbols, currentSymbol)
				currentSymbol = ""
			} else {
				invalidSyntax = true
			}

			total += ParseSymbolsValue(symbols)

			symbols = symbols[:0]
		default:
			currentSymbol += string(r)
		}
	}

	var err error
	if invalidSyntax {
		err = strconv.ErrSyntax
	}

	return total, err
}

func ParseSymbolsValue(symbols []string) float64 {
	if len(symbols) == 0 {
		return 0
	}

	max := math.Inf(-1)

	for _, symbol := range symbols {
		var value float64

		if symbol == "" {
			value = 0
		} else {
			value = GetSymbolValue(symbol)
		}

		if value > max {
			max = value
		}
	}

	return max
}

func GetSymbolValue(symbol string) float64 {
	if val, ok := KnownManaTokens[symbol]; ok {
		return val
	}

	if val, err := strconv.ParseFloat(symbol, 64); err == nil {
		return val
	}

	return 1 // Fallback: unknown symbols count as 1
}

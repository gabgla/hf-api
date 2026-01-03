package main

import (
	_ "embed"
	"encoding/gob"
	"encoding/json"
	"hf-api/src/pkg/cards"
	"hf-api/src/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const max_sides = 4 // Current maximum number of sides for any card
const sep = ";"

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing destination path argument")
	}

	destPath := os.Args[1]

	if _, err := filepath.Abs(destPath); err != nil {
		log.Fatalf("malformed path: %v", err)
	}

	var dbJSON cards.Root

	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&dbJSON); err != nil {
		log.Fatalf("decode json from stdin: %v", err)
	}

	db := normaliseDB(&dbJSON)

	file, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("failed to open file for write: %v", err)
	}

	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(db)

	if err != nil {
		log.Fatalf("failed to encode gob: %v", err)
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
			Colors:        strings.Split(c.Colors, sep),
			Sides:         parseSides(c),
			Tags:          strings.Split(c.Tags, sep),
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

		side := cards.Side{
			Cost:              utils.Coalesce(c.Cost[i], ""),
			Supertypes:        strings.Split(utils.Coalesce(c.Supertypes[i], ""), sep),
			CardTypes:         strings.Split(utils.Coalesce(c.CardTypes[i], ""), sep),
			Subtypes:          strings.Split(utils.Coalesce(c.Subtypes[i], ""), sep),
			Power:             powerInt,
			Toughness:         toughnessInt,
			Loyalty:           loyaltyInt,
			PowerOriginal:     powerStr,
			ToughnessOriginal: toughnessStr,
			LoyaltyOriginal:   loyaltyStr,
			TextBox:           utils.Coalesce(c.TextBox[i], ""),
			FlavorText:        utils.Coalesce(c.FlavorText[i], ""),

			Tags:   strings.Split(c.Tags, sep),
			Tokens: c.Tokens,
		}

		sides = append(sides, side)
	}

	return sides
}

func getStringAndNumber(field *any) (*string, *int) {
	var fieldStr string
	var fieldInt *int

	if field != nil {
		if v, ok := (*field).(string); ok {
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

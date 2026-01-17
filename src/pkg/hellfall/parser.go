package hellfall

import (
	"hf-api/src/pkg/cards"
	"hf-api/src/utils"
	"strconv"
	"strings"
)

const max_sides = 4 // Current maximum number of sides for any card
const sep = ";"

func NormaliseDB(db *Root) []cards.Card {
	result := []cards.Card{}

	for _, c := range db.Data {
		manaValueStr, manaValueNum := getStringAndNumber(c.CMC)
		card := cards.Card{
			Name:              c.Name,
			Creator:           c.Creator,
			Set:               c.Set,
			Legality:          c.ConstructedLegality,
			Rulings:           c.Rulings,
			ManaValue:         manaValueNum,
			Colors:            strings.Split(c.Colors, sep),
			Sides:             ParseSides(c),
			Tags:              strings.Split(c.Tags, sep),
			Tokens:            toDomainTokens(c.Tokens),
			ComponentOf:       c.ComponentOf,
			IsActualToken:     c.IsActualToken,
			SmallAltImage:     c.SmallAltImage,
			ManaValueOriginal: manaValueStr,
		}
		result = append(result, card)
	}

	return result
}

func ParseSides(c CardEntry) []cards.Side {
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
		if (c.Cost[i] == nil || *c.Cost[i] == "") &&
			(c.Supertypes[i] == nil || *c.Supertypes[i] == "") &&
			(c.CardTypes[i] == nil || *c.CardTypes[i] == "") &&
			(c.Subtypes[i] == nil || *c.Subtypes[i] == "") &&
			(c.Power[i] == nil || *c.Power[i] == "") &&
			(c.Toughness[i] == nil || *c.Toughness[i] == "") &&
			(c.Loyalty[i] == nil || *c.Loyalty[i] == "") &&
			(c.TextBox[i] == nil || *c.TextBox[i] == "") &&
			(c.FlavorText[i] == nil || *c.FlavorText[i] == "") {
			continue
		}

		// For numeric fields, attempt to parse integers
		// This allows for mathematical comparisons later
		manaValue, _ := cards.ParseManaValue(utils.Coalesce(c.Cost[i], ""))
		powerStr, powerNum := getStringAndNumber(c.Power[i])
		toughnessStr, toughnessNum := getStringAndNumber(c.Toughness[i])
		loyaltyStr, loyaltyNum := getStringAndNumber(c.Loyalty[i])

		side := cards.Side{
			Cost:       utils.Coalesce(c.Cost[i], ""),
			Supertypes: strings.Split(utils.Coalesce(c.Supertypes[i], ""), sep),
			CardTypes:  strings.Split(utils.Coalesce(c.CardTypes[i], ""), sep),
			Subtypes:   strings.Split(utils.Coalesce(c.Subtypes[i], ""), sep),

			Power:             powerNum,
			Toughness:         toughnessNum,
			Loyalty:           loyaltyNum,
			ManaValue:         manaValue,
			PowerOriginal:     powerStr,
			ToughnessOriginal: toughnessStr,
			LoyaltyOriginal:   loyaltyStr,
			TextBox:           utils.Coalesce(c.TextBox[i], ""),
			FlavorText:        utils.Coalesce(c.FlavorText[i], ""),

			Tags: strings.Split(c.Tags, sep),
		}

		sides = append(sides, side)
	}

	return sides
}

func getStringAndNumber(field *any) (*string, *float64) {
	var fieldStr string
	var fieldFloat *float64

	if field == nil {
		return nil, nil
	}

	if v, ok := (*field).(string); ok {
		fieldStr = v
	} else if v, ok := (*field).(float64); ok {
		fieldFloat = &v
		fieldStr = strconv.FormatFloat(v, 'f', -1, 64)
		return &fieldStr, fieldFloat
	}

	if n, err := strconv.ParseFloat(fieldStr, 64); err == nil {
		fieldFloat = &n
	}

	return &fieldStr, fieldFloat
}

func toDomainTokens(tokens []Token) []cards.Token {
	domainTokens := []cards.Token{}

	for _, t := range tokens {
		domainToken := cards.Token{
			Name:      t.Name,
			Power:     t.Power,
			Toughness: t.Toughness,
			Type:      t.Type,
			Image:     t.Image,
		}
		domainTokens = append(domainTokens, domainToken)
	}

	return domainTokens
}

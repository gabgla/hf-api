package cards

type Card struct {
	// Identifiers
	Name    string
	Creator string
	Set     string

	// Legality
	Legality []string
	Rulings  string

	// Characteristics
	ManaValue int
	Colors    []string
	Sides     []Side // most cards have 1 side; some have up to 4

	// Refs
	Tags   []string
	Tokens []Token

	// Additional fields
	ComponentOf   *string
	IsActualToken *bool
	SmallAltImage *string
}

type Side struct {
	Cost       string
	Supertypes []string
	CardTypes  []string
	Subtypes   []string
	Power      *int
	Toughness  *int
	Loyalty    *int
	TextBox    string
	FlavorText string

	Tags   []string
	Tokens []Token

	// Original fields in case numeric parsing doesn't work
	PowerOriginal     string
	ToughnessOriginal string
	LoyaltyOriginal   string
}

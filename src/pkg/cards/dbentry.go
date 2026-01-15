package cards

type Root struct {
	Data []CardEntry `json:"data"`
}

// Card represents a single card entry in the data array.
type CardEntry struct {
	Name                string    `json:"Name"`
	Image               []*string `json:"Image"` // entries can be null or ""
	Creator             string    `json:"Creator"`
	Set                 string    `json:"Set"`
	ConstructedLegality []string  `json:"Constructed"` // e.g. ["Legal"], ["Banned", ...]
	Rulings             string    `json:"Rulings"`

	CMC    *int   `json:"CMC"`      // null for lands / conspiracies
	Colors string `json:"Color(s)"` // "Red", "Blue;Black", or ""

	// Sides
	Cost       []*string `json:"Cost"`         // up to 4 entries, often with nulls
	Supertypes []*string `json:"Supertype(s)"` // e.g. "Legendary", ""
	CardTypes  []*string `json:"Card Type(s)"` // e.g. "Creature", "Land", "Sorcery"
	Subtypes   []*string `json:"Subtype(s)"`   // e.g. "Human;Noble", ""
	Power      []*any    `json:"power"`        // numbers, "*", or ""
	Toughness  []*any    `json:"toughness"`    // numbers or ""
	Loyalty    []*any    `json:"Loyalty"`      // usually "", but array is present

	TextBox    []*string `json:"Text Box"`              // up to 4 text boxes (adventures, split cards, etc.)
	FlavorText []*string `json:"Flavor Text,omitempty"` // optional, similar array structure

	Tags   string  `json:"Tags,omitempty"`   // optional tags string
	Tokens []Token `json:"tokens,omitempty"` // optional tokens array

	// Additional fields
	ComponentOf   *string `json:"Component of"`              // optional
	IsActualToken *bool   `json:"isActualToken,omitempty"`   // optional
	SmallAltImage *string `json:"small alt image,omitempty"` // optional
}

// Token represents entries in the "tokens" array for some cards.
type Token struct {
	Name      string `json:"Name"`
	Power     string `json:"Power"`
	Toughness string `json:"Toughness"`
	Type      string `json:"Type"`
	Image     string `json:"Image"`
}

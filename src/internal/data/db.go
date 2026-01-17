package data

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"hf-api/src/pkg/cards"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

const batchSize = 500

//go:embed db.gob.bin
var dbGob []byte
var db []cards.Card
var Index bleve.Index

func LoadDB() error {
	dec := gob.NewDecoder(bytes.NewReader(dbGob))
	err := dec.Decode(&db)

	if err != nil {
		return err
	}

	start := time.Now()

	index, _ := bleve.NewMemOnly(buildMapping())

	batch := index.NewBatch()
	for i, c := range db {
		batch.Index(c.Name, &c)

		if (i+1)%500 == 0 {
			_ = index.Batch(batch)
			batch = index.NewBatch()
		}
	}
	if batch.Size() > 0 {
		_ = index.Batch(batch)
	}

	Index = index

	end := time.Now()
	elapsed := end.Sub(start)
	println("Indexed", len(db), "cards in", elapsed.String())

	return nil
}

func buildMapping() *mapping.IndexMappingImpl {
	m := bleve.NewIndexMapping()
	m.DefaultAnalyzer = "en"

	doc := bleve.NewDocumentMapping()

	name := bleve.NewTextFieldMapping()
	name.Store = true // useful to display result titles without DB lookup
	doc.AddFieldMappingsAt("name", name)

	// Exact-match/filter fields
	kw := bleve.NewKeywordFieldMapping()
	doc.AddFieldMappingsAt("set", kw)
	doc.AddFieldMappingsAt("creator", kw)
	doc.AddFieldMappingsAt("colors", kw)
	doc.AddFieldMappingsAt("tags", kw)
	doc.AddFieldMappingsAt("legality", kw)

	// Numeric range filters / sorting
	num := bleve.NewNumericFieldMapping()
	doc.AddFieldMappingsAt("mv", num)

	m.DefaultMapping = doc
	return m
}

// func translateScryfallToBleveQS(in string) string {
// 	fieldMap := map[string]string{
// 		"o": "oracle", "oracle": "oracle",
// 		"t": "type_line", "type": "type_line",
// 		"set": "set",
// 		"c":   "colors", "color": "colors",
// 		"mv": "mv", "cmc": "mv",
// 		"tag":   "tags",
// 		"legal": "legality",
// 		"name":  "name",
// 	}

// 	toks := splitTokens(in) // handles quotes and parentheses minimally
// 	var out []string

// 	for _, tok := range toks {
// 		if tok == "or" || tok == "OR" {
// 			out = append(out, "OR")
// 			continue
// 		}

// 		neg := strings.HasPrefix(tok, "-")
// 		if neg {
// 			tok = strings.TrimPrefix(tok, "-")
// 		}

// 		// fielded forms: key:value  OR  key<=3  OR  key:<=3
// 		key, op, val, ok := parseField(tok)
// 		if ok {
// 			key = strings.ToLower(key)
// 			if mapped, exists := fieldMap[key]; exists {
// 				key = mapped
// 			}

// 			// Expand c:uw => +colors:U +colors:W
// 			if key == "colors" && op == ":" && isColorLetters(val) {
// 				for _, r := range strings.ToUpper(val) {
// 					clause := fmt.Sprintf(`%scolors:%c`, prefix(neg), r)
// 					// Scryfall implicit AND -> force each color required unless negated
// 					if !neg {
// 						clause = "+" + clause
// 					}
// 					out = append(out, clause)
// 				}
// 				continue
// 			}

// 			// Normal field clause
// 			// Make it required by default (Scryfall semantics), unless it’s inside an OR group you build explicitly.
// 			clause := fmt.Sprintf(`%s%s:%s%s`, prefix(neg), key, op, val)
// 			if !neg {
// 				clause = "+" + clause
// 			}
// 			out = append(out, clause)
// 			continue
// 		}

// 		// Bare term / phrase / regex: require it to mimic Scryfall default AND.
// 		if neg {
// 			out = append(out, "-"+tok)
// 		} else {
// 			out = append(out, "+"+tok)
// 		}
// 	}

// 	return strings.Join(out, " ")
// }

// func prefix(neg bool) string {
// 	if neg {
// 		return "-"
// 	}
// 	return ""
// }

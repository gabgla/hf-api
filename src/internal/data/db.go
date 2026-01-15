package data

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"hf-api/src/pkg/cards"

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

	index, _ := bleve.NewMemOnly(buildMapping())

	batch := index.NewBatch()
	for i, c := range db {
		doc := toDoc(c.Name, c)
		batch.Index(doc.ID, doc)

		if (i+1)%500 == 0 {
			_ = index.Batch(batch)
			batch = index.NewBatch()
		}
	}
	if batch.Size() > 0 {
		_ = index.Batch(batch)
	}

	Index = index

	return nil
}

type CardDoc struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Creator   string   `json:"creator"`
	Set       string   `json:"set"`
	ManaValue float64  `json:"mv"`
	Colors    []string `json:"colors"`
	Legality  []string `json:"legality"`
	Tags      []string `json:"tags"`
}

func toDoc(id string, c cards.Card) CardDoc {
	return CardDoc{
		ID:        id,
		Name:      c.Name,
		Creator:   c.Creator,
		Set:       c.Set,
		ManaValue: float64(c.ManaValue),
		Colors:    c.Colors,
		Legality:  c.Legality,
		Tags:      c.Tags,
	}
}

func buildMapping() *mapping.IndexMappingImpl {
	m := bleve.NewIndexMapping()
	m.DefaultAnalyzer = "en"

	doc := bleve.NewDocumentMapping()

	// Full-text
	ft := bleve.NewTextFieldMapping()
	ft.Store = false
	doc.AddFieldMappingsAt("full_text", ft)

	// Also index name separately (so you can boost it)
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
var DBGob []byte

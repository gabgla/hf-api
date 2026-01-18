package main

import (
	_ "embed"
	"encoding/gob"
	"encoding/json"
	"hf-api/src/pkg/cards"
	"hf-api/src/pkg/hellfall"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

const batchSize = 500
const dbGobFilename = "db.gob.bin"
const indexName = "index.bleve"

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing source path argument")
	}

	if len(os.Args) < 3 {
		log.Fatalf("missing destination path argument")
	}

	sourcePath := os.Args[1]
	destPath := os.Args[2]

	if _, err := filepath.Abs(sourcePath); err != nil {
		log.Fatalf("malformed path: %v", err)
	}

	if _, err := filepath.Abs(destPath); err != nil {
		log.Fatalf("malformed path: %v", err)
	}

	var dbJSON hellfall.Root

	sourceFile, err := os.Open(sourcePath)

	if err != nil {
		log.Fatalf("failed to open source file: %v", err)
	}

	defer sourceFile.Close()

	dec := json.NewDecoder(sourceFile)
	if err := dec.Decode(&dbJSON); err != nil {
		log.Fatalf("decode json from stdin: %v", err)
	}

	db := hellfall.NormaliseDB(&dbJSON)

	if err := writeDB(filepath.Join(destPath, dbGobFilename), db); err != nil {
		log.Fatalf("failed to encode gob: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	if err := generateIndex(filepath.Join(wd, indexName), db); err != nil {
		log.Fatalf("failed to generate index: %v", err)
	}
}

func writeDB(destPath string, db []cards.Card) error {
	destFile, err := os.Create(destPath)

	if err != nil {
		log.Fatalf("failed to open file for write: %v", err)
	}

	defer destFile.Close()

	enc := gob.NewEncoder(destFile)
	return enc.Encode(db)
}

func generateIndex(indexPath string, db []cards.Card) error {
	os.RemoveAll(indexPath)

	mapping := buildMapping()
	index, err := bleve.New(indexPath, mapping)

	if err != nil {
		return err
	}

	batch := index.NewBatch()
	for i, c := range db {
		docID := strconv.Itoa(i)
		batch.Index(docID, &c)

		if (i+1)%batchSize == 0 {
			_ = index.Batch(batch)
			batch = index.NewBatch()
		}
	}
	if batch.Size() > 0 {
		_ = index.Batch(batch)
	}

	if err := index.Close(); err != nil {
		return err
	}

	return nil
}

func buildMapping() *mapping.IndexMappingImpl {
	m := bleve.NewIndexMapping()
	m.DefaultAnalyzer = "en"

	doc := bleve.NewDocumentMapping()

	subDoc := bleve.NewDocumentMapping()
	doc.AddSubDocumentMapping("sides", subDoc)

	// Text fields
	text := bleve.NewTextFieldMapping()
	doc.AddFieldMappingsAt("name", text)

	// Keyword fields
	keyword := bleve.NewKeywordFieldMapping()
	doc.AddFieldMappingsAt("creator", keyword)
	doc.AddFieldMappingsAt("set", keyword)
	doc.AddFieldMappingsAt("legality", keyword)
	doc.AddFieldMappingsAt("colors", keyword)
	doc.AddFieldMappingsAt("tags", keyword)
	doc.AddFieldMappingsAt("mv_original", keyword)

	// Numeric fields
	num := bleve.NewNumericFieldMapping()
	doc.AddFieldMappingsAt("mv", num)

	// Text fields for sides
	subDoc.AddFieldMappingsAt("cost", text)
	subDoc.AddFieldMappingsAt("supertypes", text)
	subDoc.AddFieldMappingsAt("card_types", text)
	subDoc.AddFieldMappingsAt("subtypes", text)
	subDoc.AddFieldMappingsAt("textbox", text)
	subDoc.AddFieldMappingsAt("flavor_text", text)

	// Keyword fields for sides
	subDoc.AddFieldMappingsAt("mv_original", keyword)
	subDoc.AddFieldMappingsAt("power_original", keyword)
	subDoc.AddFieldMappingsAt("toughness_original", keyword)
	subDoc.AddFieldMappingsAt("loyalty_original", keyword)

	// Numeric fields for sides
	subDoc.AddFieldMappingsAt("mv", num)
	subDoc.AddFieldMappingsAt("power", num)
	subDoc.AddFieldMappingsAt("toughness", num)
	subDoc.AddFieldMappingsAt("loyalty", num)

	m.DefaultMapping = doc

	return m
}

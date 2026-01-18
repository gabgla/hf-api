package data

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"fmt"
	"hf-api/src/pkg/cards"
	"time"

	"github.com/blevesearch/bleve/v2"
)

const batchSize = 500
const indexName = "index.bleve"

//go:embed db.gob.bin
var dbGob []byte

var DB []cards.Card
var Index bleve.Index

func LoadDB() error {
	start := time.Now()

	dec := gob.NewDecoder(bytes.NewReader(dbGob))
	err := dec.Decode(&DB)

	if err != nil {
		return err
	}

	end := time.Now()
	elapsed := end.Sub(start)
	println("DB loaded in", elapsed.Milliseconds(), "ms")
	fmt.Printf("loaded %d items in %dms\n", len(DB), elapsed.Milliseconds())

	start = time.Now()
	Index, err = bleve.Open(indexName)
	end = time.Now()
	elapsed = end.Sub(start)

	if err != nil {
		return err
	}

	println("Index loaded in", elapsed.Milliseconds(), "ms")

	return nil
}

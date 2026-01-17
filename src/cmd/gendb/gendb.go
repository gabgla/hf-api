package main

import (
	_ "embed"
	"encoding/gob"
	"encoding/json"
	"hf-api/src/pkg/hellfall"
	"log"
	"os"
	"path/filepath"
)

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

	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("failed to open file for write: %v", err)
	}

	defer destFile.Close()

	enc := gob.NewEncoder(destFile)
	err = enc.Encode(db)

	if err != nil {
		log.Fatalf("failed to encode gob: %v", err)
	}
}

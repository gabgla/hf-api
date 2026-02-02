DATABASE_URL="https://raw.githubusercontent.com/bones-bones/hellfall/main/src/data/Hellscube-Database.json"
DATABASE_JSON_FILENAME=database.json
BUILD_FLAGS := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
LDFLAGS := -trimpath -ldflags="-s -w"

download-db:
	curl -o "./${DATABASE_JSON_FILENAME}" "${DATABASE_URL}"

generate-db:
	go run src/cmd/gendb/gendb.go "${DATABASE_JSON_FILENAME}" "src/internal/data"

generate-token-aliases:
	go run src/cmd/codegens/codegens.go -- "src/internal/app/api/handlers_tokens.go"

setup: download-db generate-db

build-hfapi: # Don't change this name, it's used by AWS SAM
	$(BUILD_FLAGS) go build $(LDFLAGS) -o $(ARTIFACTS_DIR)/bootstrap ./src/cmd/lambda
	cp -R index.bleve $(ARTIFACTS_DIR)/

build-for-lambda: setup clean
	mkdir -p build/lambda
	$(BUILD_FLAGS) go build $(LDFLAGS) -o build/lambda/bootstrap ./src/cmd/lambda
	cp -R index.bleve build/lambda/

run-http:
	go run ./src/cmd/httpserver/httpserver.go

run-lambda:
	sam build --cached && sam local start-api --warm-containers EAGER

run: run-http

test:
	go test ./...

clean:
	rm -rf build/

.PHONY: download-db generate-db generate-token-aliases setup build-hfapi build-for-lambda run-http run-lambda run test clean

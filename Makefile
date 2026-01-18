DATABASE_URL="https://raw.githubusercontent.com/bones-bones/hellfall/main/src/data/Hellscube-Database.json"
DATABASE_JSON_FILENAME=database.json

download-db:
	curl -o "./${DATABASE_JSON_FILENAME}" "${DATABASE_URL}"

generate-db:
	go run src/cmd/gendb/gendb.go "${DATABASE_JSON_FILENAME}" "src/internal/data"

generate-token-aliases:
	go run src/cmd/codegens/codegens.go -- "src/internal/app/api/handlers_tokens.go"

setup: download-db generate-db

build-hfapi:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(ARTIFACTS_DIR)/bootstrap ./src/cmd/lambda
	cp -R index.bleve $(ARTIFACTS_DIR)/

run-http:
	go run ./src/cmd/httpserver/httpserver.go

run-lambda:
	sam build --cached && sam local start-api --warm-containers EAGER

run: run-http

.PHONY: download-db generate-db generate-token-aliases setup build-hfapi-lambda run-http run-lambda run

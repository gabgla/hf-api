DATABASE_URL="https://raw.githubusercontent.com/bones-bones/hellfall/main/src/data/Hellscube-Database.json"
DATABASE_JSON_FILENAME=database.json
DATABASE_GOB_FILENAME=db.gob.bin

download-db:
	curl -o "./${DATABASE_JSON_FILENAME}" "${DATABASE_URL}"

generate-db:
	go run src/cmd/gendb/gendb.go "src/internal/data/${DATABASE_GOB_FILENAME}" < "${DATABASE_JSON_FILENAME}"

generate-token-aliases:
	go run src/cmd/codegens/codegens.go -- "src/internal/app/api/handlers_tokens.go"

setup: download-db generate-db

run-http:
	go run ./src/cmd/httpserver/httpserver.go

run-lambda:
	go run ./src/cmd/lambda/lambda.go

run: run-http

.PHONY: download-db generate-db generate-token-aliases setup run-http run-lambda run

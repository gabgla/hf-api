DATABASE_URL="https://raw.githubusercontent.com/bones-bones/hellfall/main/src/data/Hellscube-Database.json"

.PHONY: download-db
download-db:
	curl -o ./database.json ${DATABASE_URL}

.PHONY: run
run:
	LOCAL=1 go run ./src/cmd/main.go

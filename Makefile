.PHONY: \
	backend-run backend-tidy backend-test backend-fmt backend-vet \
	mysql-up mysql-down mysql-reset mysql-logs mysql-ps \
	health \
	contracts-build contracts-test contracts-clean \
	git-status

backend-run:
	cd backend && go run ./cmd/server

backend-tidy:
	cd backend && go mod tidy

backend-test:
	cd backend && go test ./...

backend-fmt:
	cd backend && gofmt -w .

backend-vet:
	cd backend && go vet ./...

mysql-up:
	cd backend && docker compose up -d

mysql-down:
	cd backend && docker compose down

mysql-reset:
	cd backend && docker compose down -v && docker compose up -d

mysql-logs:
	cd backend && docker compose logs -f mysql

mysql-ps:
	cd backend && docker compose ps

health:
	curl http://localhost:8080/health

contracts-build:
	cd contracts && forge build

contracts-test:
	cd contracts && forge test

contracts-clean:
	cd contracts && forge clean

git-status:
	git status
.PHONY: backend-run backend-tidy mysql-up mysql-down mysql-reset mysql-logs mysql-ps health git-status

backend-run:
	cd backend && go run ./cmd/server

backend-tidy:
	cd backend && go mod tidy

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

git-status:
	git status
build:
	go build -o main .

test:
	go test ./...

infra-up: # including test db
	docker compose up -d --build
	sleep 5
	migrate -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path migrations/ up

infra-down:
	docker compose down --remove-orphans

docker-build:
	docker build -t itspay:latest .

run:
	go run main.go

run-docker-compose: infra-up
	docker compose run --rm -P app ./app

lint:
	golangci-lint run -v

generate:
	go generate ./...
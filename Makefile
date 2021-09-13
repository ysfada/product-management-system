run:
	go run . # -debug

build:
	go build -o ./bin/ .

test:
	go test -coverprofile cp.out ./... && go tool cover -html=cp.out -o coverage.html

up:
	docker-compose up -d

down:
	docker-compose down

migrate-new:
	migrate create -ext sql -dir database/migrations $(name)

migrate-up:
	migrate -database postgres://username:password@localhost:5432/dbname?sslmode=disable -path database/migrations up

migrate-down:
	migrate -database postgres://username:password@localhost:5432/dbname?sslmode=disable -path database/migrations down

init:
	swag init

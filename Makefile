.PHONY: up init build-billing build-queue test migration migrate-up swagger

up:
	docker-compose up -d --build

init:
	go mod init github.com/iseroukhov/brave-new-billing

build-billing:
	go mod download && go mod verify && go mod tidy && go mod vendor && go build -mod=vendor -o ./bin/billing ./cmd/billing

build-queue:
	go mod download && go mod verify && go mod tidy && go mod vendor && go build -mod=vendor -o ./bin/queue ./cmd/queue

test:
	go test -v -mod=vendor -coverpkg=./... ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

migration:
	migrate create -ext sql -dir database/migrations $(t)

migrate-up:
	migrate -path "database/migrations" -database "mysql://root:root@tcp(localhost:33060)/billing?charset=utf8&interpolateParams=true" up

swagger:
	swagger serve api/swagger-spec/api.json
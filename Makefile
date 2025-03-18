.PHONY: proto build-auth build-chat build-gateway build run test lint docker clean migrate-up migrate-down swag

# Генерация Protocol Buffers
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/*.proto proto/chat/*.proto

# Сборка сервисов
build-auth:
	cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go

build-chat:
	cd services/chat-service && go build -o ../../bin/chat-service ./cmd/main.go

build-gateway:
	cd services/gateway-service && go build -o ../../bin/gateway-service ./cmd/main.go

build: build-auth build-chat build-gateway

# Запуск Docker Compose
run:
	docker-compose up -d

# Запуск без Docker (локально)
run-local: build
	bin/auth-service & bin/chat-service & bin/gateway-service

# Тестирование
test:
	go test -v ./services/auth-service/... ./services/chat-service/... ./services/gateway-service/... ./proto/...

# Тестирование с race detector
test-race:
	go test -race -v ./services/auth-service/... ./services/chat-service/... ./services/gateway-service/... ./proto/...

# Линтер
lint:
	golangci-lint run ./...

# Сборка Docker образов
docker:
	docker build -t gochat/auth-service -f services/auth-service/Dockerfile .
	docker build -t gochat/chat-service -f services/chat-service/Dockerfile .
	docker build -t gochat/gateway-service -f services/gateway-service/Dockerfile .

# Генерация Swagger документации
swag:
	cd services/gateway-service && swag init -g cmd/main.go

# Миграции
migrate-up:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/chatapp?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/chatapp?sslmode=disable" down

# Очистка
clean:
	rm -rf bin/
	docker-compose down 
.PHONY: proto build-auth build-chat build-gateway build run test lint docker clean migrate-up migrate-down swag setup-workspace

# Настройка Go workspace
setup-workspace:
	@echo "Пропускаем go work sync из-за конфликтов"

# Генерация Protocol Buffers
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/*.proto proto/chat/*.proto

# Сборка сервисов
build-auth: setup-workspace
	cd services/auth-service && GOWORK=off go build -o ../../bin/auth-service ./cmd/main.go

build-chat: setup-workspace
	cd services/chat-service && GOWORK=off go build -o ../../bin/chat-service ./cmd/main.go

build-gateway: setup-workspace
	cd services/gateway-service && GOWORK=off go build -o ../../bin/gateway-service ./cmd/main.go

build: setup-workspace build-auth build-chat build-gateway

# Запуск Docker Compose
run:
	docker-compose up -d

# Запуск без Docker (локально)
run-local: build
	bin/auth-service & bin/chat-service & bin/gateway-service

# Тестирование
test: setup-workspace
	GOWORK=off go test -v ./pkg/... ./services/auth-service/... ./services/chat-service/... ./services/gateway-service/... ./proto/...

# Тестирование с race detector
test-race: setup-workspace
	GOWORK=off go test -race -v ./pkg/... ./services/auth-service/... ./services/chat-service/... ./services/gateway-service/... ./proto/...

# Линтер
lint: setup-workspace
	GOWORK=off golangci-lint run ./pkg/...
	GOWORK=off golangci-lint run ./proto/...
	GOWORK=off golangci-lint run ./services/...

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
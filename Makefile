.PHONY: proto build-auth build-chat build-gateway build run test lint docker clean migrate-up migrate-down swag setup-workspace

# Настройка Go workspace
setup-workspace:
	go work sync

# Генерация Protocol Buffers
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/*.proto proto/chat/*.proto

# Сборка сервисов
build-auth: setup-workspace
	cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go

build-chat: setup-workspace
	cd services/chat-service && go build -o ../../bin/chat-service ./cmd/main.go

build-gateway: setup-workspace
	cd services/gateway-service && go build -o ../../bin/gateway-service ./cmd/main.go

build: setup-workspace build-auth build-chat build-gateway

# Запуск Docker Compose
run:
	docker-compose up -d

# Запуск без Docker (локально)
run-local: build
	bin/auth-service & bin/chat-service & bin/gateway-service

# Тестирование
test: setup-workspace
	go test -v ./pkg/... ./services/auth-service/... ./services/chat-service/... ./services/gateway-service/... ./proto/...

# Тестирование с race detector
test-race: setup-workspace
	go test -race -v ./pkg/... ./services/auth-service/... ./services/chat-service/... ./services/gateway-service/... ./proto/...

# Линтер
lint: setup-workspace
	golangci-lint run ./pkg/...
	golangci-lint run ./proto/...
	golangci-lint run ./services/...

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
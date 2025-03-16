.PHONY: proto build-auth build-chat build-gateway build run clean

# Генерация Protocol Buffers
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/*.proto proto/chat/*.proto

# Сборка сервисов
build-auth:
	cd services/auth-service && go build -o ../../bin/auth-service

build-chat:
	cd services/chat-service && go build -o ../../bin/chat-service

build-gateway:
	cd services/gateway-service && go build -o ../../bin/gateway-service

build: build-auth build-chat build-gateway

# Запуск с помощью Docker Compose
run:
	docker-compose up -d

# Очистка
clean:
	rm -rf bin/
	docker-compose down 
.PHONY: proto build run test clean

# Генерация gRPC кода из .proto файлов
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Сборка всех сервисов
build:
	cd services/auth-service && go build -o ../../bin/auth-service ./cmd
	cd services/chat-service && go build -o ../../bin/chat-service ./cmd
	cd services/gateway-service && go build -o ../../bin/gateway-service ./cmd

# Запуск с помощью docker-compose
run:
	docker-compose up -d

# Тестирование
test:
	go test ./...

# Очистка
clean:
	rm -rf bin/*
	docker-compose down -v 
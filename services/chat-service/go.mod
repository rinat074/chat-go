module github.com/rinat074/chat-go/services/chat-service

go 1.21

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/jackc/pgx/v5 v5.4.3
	github.com/rinat074/chat-go/proto/chat v0.0.0
	google.golang.org/grpc v1.58.2
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230920204549-e6e6cdab5c13 // indirect
)

replace (
	github.com/rinat074/chat-go/pkg => ../../pkg
	github.com/rinat074/chat-go/proto/chat => ../../proto/chat
)

module github.com/rinat074/chat-go

go 1.21

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/gorilla/websocket v1.5.0
	github.com/redis/go-redis/v9 v9.1.0
	github.com/rinat074/chat-go/pkg v0.0.0
	github.com/rinat074/chat-go/proto/auth v0.0.0
	github.com/rinat074/chat-go/proto/chat v0.0.0
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.2
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.4
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/rinat074/chat-go/pkg => ./pkg
	github.com/rinat074/chat-go/proto/auth => ./proto/auth
	github.com/rinat074/chat-go/proto/chat => ./proto/chat
)

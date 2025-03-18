module github.com/rinat074/chat-go/services/auth-service

go 1.21

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/jackc/pgx/v5 v5.4.3
	github.com/rinat074/chat-go/pkg v0.0.0
	github.com/rinat074/chat-go/proto/auth v0.0.0
	golang.org/x/crypto v0.14.0
	google.golang.org/grpc v1.58.2
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230920204549-e6e6cdab5c13 // indirect
)

replace (
	github.com/rinat074/chat-go/pkg => ../../pkg
	github.com/rinat074/chat-go/proto/auth => ../../proto/auth
)

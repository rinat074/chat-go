module chat-app/services/auth-service

go 1.23.0

toolchain go1.23.7

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gochat/proto/auth v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.7.2
	golang.org/x/crypto v0.36.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
)

replace github.com/gochat/proto/auth => ../../proto/auth

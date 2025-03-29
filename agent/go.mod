module reconmap/agent

go 1.23.0

toolchain go1.24.1

require (
	github.com/creack/pty v1.1.24
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/reconmap/shared-lib v0.0.0-20220910165932-7d018d9111fc
	github.com/robfig/cron v1.2.0
	go.uber.org/zap v1.27.0
)

require github.com/golang-jwt/jwt/v5 v5.2.2 // indirect

require (
	github.com/Nerzal/gocloak/v13 v13.9.0
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.38.0 // indirect
)

replace github.com/reconmap/shared-lib => ../shared-lib

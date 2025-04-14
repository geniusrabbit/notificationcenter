module github.com/geniusrabbit/notificationcenter/v2/redis

go 1.23.0

toolchain go1.24.2

require (
	github.com/alicebob/miniredis/v2 v2.34.0
	github.com/demdxx/gocast/v2 v2.9.0
	github.com/geniusrabbit/notificationcenter/v2 v2.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.7.3
	github.com/stretchr/testify v1.10.0
	go.uber.org/multierr v1.11.0
)

require (
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/geniusrabbit/notificationcenter/v2 => ../

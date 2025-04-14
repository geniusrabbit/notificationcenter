module github.com/geniusrabbit/notificationcenter/v2/nats

go 1.23.0

toolchain go1.24.2

require (
	github.com/geniusrabbit/notificationcenter/v2 v2.0.0-00010101000000-000000000000
	github.com/nats-io/nats.go v1.41.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/nkeys v0.4.10 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/geniusrabbit/notificationcenter/v2 => ../

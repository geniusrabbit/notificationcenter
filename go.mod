module github.com/geniusrabbit/notificationcenter/v2

go 1.24.0

require (
	github.com/IBM/sarama v1.46.1 // for kafka
	github.com/alicebob/miniredis/v2 v2.35.0
	github.com/allegro/bigcache v1.2.1 // for wrappers/once/bigcache
	github.com/demdxx/gocast/v2 v2.10.2
	github.com/demdxx/rpool/v2 v2.0.1 // for wrappers/concurrency
	github.com/golang/mock v1.6.0
	github.com/lib/pq v1.10.9 // for pg
	github.com/nats-io/nats.go v1.45.0 // for nats, natstream
	github.com/nats-io/stan.go v0.10.4 // for nats, natstream
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.14.0 // for redis, wrappers/once/redis
	github.com/stretchr/testify v1.11.1
	go.uber.org/multierr v1.11.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/nats-server/v2 v2.11.9 // indirect
	github.com/nats-io/nats-streaming-server v0.25.6 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

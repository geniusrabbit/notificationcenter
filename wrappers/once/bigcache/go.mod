module github.com/geniusrabbit/notificationcenter/v2/wrappers/once/bigcache

go 1.23.0

toolchain go1.24.2

require (
	github.com/allegro/bigcache v1.2.1
	github.com/geniusrabbit/notificationcenter/v2 v2.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/geniusrabbit/notificationcenter/v2 => ../../../

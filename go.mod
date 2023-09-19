module polar-sky

go 1.19

require (
	github.com/docker/go-connections v0.4.0
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/google/uuid v1.3.1
	go.uber.org/zap v1.26.0
)

require (
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace github.com/uber-go/zap v1.26.0 => go.uber.org/zap v1.26.0

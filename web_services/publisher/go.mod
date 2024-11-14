module github.com/saiddis/practicing_go/web_services/publisher

go 1.23.0

replace github.com/saiddis/practicing_go/web_services/publisher/pkg/server => ./pkg/server

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/cors v1.2.1
	github.com/nats-io/nats.go v1.37.0
)

require (
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

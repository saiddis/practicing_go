package internal

import "github.com/nats-io/nats.go"

type ApiConfig struct {
	NC *nats.Conn
}

func NewApi(nc *nats.Conn) *ApiConfig {
	return &ApiConfig{
		NC: nc,
	}
}

package server

import (
	"errors"
	"net/http"
)

type options struct {
	port     int
	handler  http.Handler
	natsAddr string
}

type Option func(opts *options) error

func WithPort(port int) Option {
	return func(opts *options) error {
		if port <= 0 {
			return errors.New("port parameter has to be more than 0")
		}
		opts.port = port
		return nil
	}
}

func WithNatsAddr(addr string) Option {
	return func(opts *options) error {
		opts.natsAddr = addr
		return nil
	}
}

package server

import (
	"errors"
	"net/http"
)

type options struct {
	port    int
	handler http.Handler
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

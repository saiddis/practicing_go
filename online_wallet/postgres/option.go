package postgres

type parameters struct {
	sslMode  string
	timeZone string
}

type options struct {
	url      string
	user     string
	host     string
	password string
	port     int
	params   map[string]string
}

type Option func(opts *options)

func WithURL(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}

func WithHost(host string) Option {
	return func(opts *options) {
		opts.host = host
	}
}

func WithUser(usr string) Option {
	return func(opts *options) {
		opts.user = usr
	}
}

func WithPassword(pswd string) Option {
	return func(opts *options) {
		opts.password = pswd
	}
}

func WithPort(port int) Option {
	return func(opts *options) {
		opts.port = port
	}
}

func WithSSL(ssl string) Option {
	return func(opts *options) {
		opts.params["sslmode"] = ssl
	}
}

func WithTimeZone(tz string) Option {
	return func(opts *options) {
		opts.params["TimeZone"] = tz
	}
}

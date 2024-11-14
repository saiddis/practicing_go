package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/nats-io/nats.go"
	"github.com/saiddis/practicing_go/web_services/service1/internal"
)

type Server struct {
	engine    *http.Server
	apiConfig *internal.ApiConfig
}

func New(addr string, opts ...Option) (*Server, error) {
	options := options{
		port:     8080,
		natsAddr: nats.DefaultURL,
	}
	var err error
	for _, opt := range opts {
		err = opt(&options)
		if err != nil {
			return nil, err
		}
	}
	server := &http.Server{
		Addr: ":" + strconv.Itoa(options.port),
	}
	nc, err := nats.Connect(options.natsAddr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS server: %v", err)
	}
	apiConfig := internal.NewApi(nc)
	return &Server{engine: server, apiConfig: apiConfig}, nil
}

func (s *Server) Run() error {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Post("/greet", s.apiConfig.Greet)

	s.engine.Handler = router

	log.Printf("Listening on %s", s.engine.Addr)
	err := s.engine.ListenAndServe()
	if err != nil {
		return fmt.Errorf("Couln't run the server: %v", err)
	}
	return nil
}

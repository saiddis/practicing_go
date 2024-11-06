package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/saiddis/practicing_go/online_wallet/controllers"
	"github.com/saiddis/practicing_go/online_wallet/domain"
)

type Server struct {
	userUsecase *controllers.UserUsecase
	engine      *http.Server
}

func New(addr string, userUsecase *controllers.UserUsecase, opts ...Option) (*Server, error) {
	options := options{
		port: 8080,
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

	return &Server{userUsecase: userUsecase, engine: server}, nil
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

	router.Route("/user", func(r chi.Router) {
		// User routes
		r.Get("/", s.userUsecase.GetUser)
		r.Post("/", s.userUsecase.CreateUser)
		r.Patch("/", s.userUsecase.UpdateUser)
		r.Delete("/", s.userUsecase.DeleteUser)

		// Wallet routes
		r.Post("/credit", s.userUsecase.AddUpToBalance)
		r.Post("/transfer", s.userUsecase.Transfer)
	})

	s.engine.Handler = router

	log.Printf("Running on port%s", s.engine.Addr)
	err := s.engine.ListenAndServe()
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "Couln't run the server: %v", err)
	}
	return nil
}

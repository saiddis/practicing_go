package main

import (
	"log"

	"github.com/saiddis/practicing_go/online_wallet/controllers"
	"github.com/saiddis/practicing_go/online_wallet/postgres"
	"github.com/saiddis/practicing_go/online_wallet/repository"
	"github.com/saiddis/practicing_go/online_wallet/server"
)

func main() {
	db, err := postgres.New("wallet",
		postgres.WithUser("saiddis"),
		postgres.WithPassword("__1dIslo_"),
		postgres.WithSSL("disable"),
		postgres.WithTimeZone("Asia/Dushanbe"),
	)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %v", err)
	}
	userRepository := repository.NewUserRepository(db)
	userUsecase := controllers.NewUserUsecase(userRepository)
	server, err := server.New("localhost", userUsecase, server.WithPort(8080))
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

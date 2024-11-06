package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/saiddis/practicing_go/online_wallet/domain"
)

type UserUsecase struct {
	repository domain.UserService
}

func NewUserUsecase(repo domain.UserService) *UserUsecase {
	return &UserUsecase{
		repository: repo,
	}
}

func (u *UserUsecase) CreateUser(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		Name    string  `json:"name"`
		Email   string  `json:"email"`
		Balance float32 `json:"balance"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user := domain.User{
		Name:      parameters.Name,
		Email:     parameters.Email,
		Balance:   parameters.Balance,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	log.Printf("user: %+v", user)

	ctx := r.Context()
	err = u.repository.CreateUser(ctx, user)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJSON(w, 201, user)
}

func (u *UserUsecase) GetUser(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		ID int `json:"id"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ctx := r.Context()
	user, err := u.repository.FindUserByID(ctx, parameters.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving user: %v", err))
		return
	}
	respondWithJSON(w, 200, user)
}

func (u *UserUsecase) AddUpToBalance(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		ID     int     `json:"id"`
		Amount float32 `json:"amount"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ctx := r.Context()
	err = u.repository.Credit(ctx, parameters.ID, parameters.Amount)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error adding up to balance: %v", err))
		return
	}
	respondWithJSON(w, 200, fmt.Sprintf("%f added to balance of user with id %d", parameters.Amount, parameters.ID))
}

func (u *UserUsecase) Transfer(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		From   int     `json:"from"`
		To     int     `json:"to"`
		Amount float32 `json:"amount"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ctx := r.Context()
	err = u.repository.Transfer(ctx, parameters.From, parameters.To, parameters.Amount)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error transfering money: %v", err))
		return
	}
	respondWithJSON(w, 200, fmt.Sprintf("%f transfered from user with id %d to user with id %d", parameters.Amount, parameters.From, parameters.To))
}

func (u *UserUsecase) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ctx := r.Context()
	err = u.repository.UpdateUser(ctx, parameters.ID, domain.UserUpdate{
		Name:  parameters.Name,
		Email: parameters.Email,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error updating user: %v", err))
		return
	}
	respondWithJSON(w, 200, parameters)
}

func (u *UserUsecase) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		ID int `json:"id"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ctx := r.Context()
	err = u.repository.DeleteUser(ctx, parameters.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting user: %v", err))
		return
	}
	respondWithJSON(w, 200, fmt.Sprintf("user with id %d was successfully deleted", parameters.ID))
}

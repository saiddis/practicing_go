package domain

import (
	"context"
	"time"
)

const (
	UnauthorizedBalanceLimit = 10_000
	AuthorizedBalanceLimit   = 100_000
)

// User represents a user on the system.
type User struct {
	ID int `json:"id"`

	// Users prefered name and email.
	Name  string `json:"name"`
	Email string `json:"email"`

	// Timestamps for user creation and last update.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// User's balance
	Balance float32 `json:"balance"`
}

// Validate returns an error if the user has invalid fields.
func (u User) Validate() error {
	if u.Name == "" {
		return Errorf(EINVALID, "User name required.")
	}

	return nil
}

// UserService represents a service for managing users.
type UserService interface {
	// Retrieves a user by ID along with their associated object.
	FindUserByID(ctx context.Context, id int) (*User, error)

	// Creates a new user.
	CreateUser(ctx context.Context, user User) error

	// Updates a user object.
	UpdateUser(ctx context.Context, id int, upd UserUpdate) error

	// Permanently a user and all owned repos.
	DeleteUser(ctx context.Context, id int) error

	// Transfers money from user's wallet to another.
	Transfer(ctx context.Context, dstID, srcID int, amount float32) error

	// Adds up money to the wallet. Returns EREACHEDLIMIT if a user is not
	// authorized and its balance + credit is more than  the limit of 10 000.
	// In case the user is authorized the limit is 100 000.
	Credit(ctx context.Context, id int, amount float32) error

	// Withdraws money from user balance with the given id. Returns ENOTENOUGH if
	// the given amount is greater than balance.
	Withdraw(ctx context.Context, id int, amount float32) error
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

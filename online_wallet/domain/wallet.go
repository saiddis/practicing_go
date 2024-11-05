package domain

import "context"

// Wallet represents an online wallet.
type Wallet struct {
	Balance float64 `json:"balance"`
}

// WalletService represents a service for managing wallets.
type WalletService interface {
	// Transfers money from user's wallet to another.
	Transfer(ctx context.Context, dst User, amount float64) error

	// Adds up money to the wallet. Returns EREACHEDLIMIT if a user is not
	// authorized and its balance + credit is more than  the limit of 10 000.
	// In case the user is authorized the limit is 100 000.
	Credit(ctx context.Context, amount float64) error
}

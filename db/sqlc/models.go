// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"time"
)

type Account struct {
	ID        int32     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int32     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int32 `json:"id"`
	AccountID int32 `json:"account_id"`
	// can be negative or positive
	Amount    int32     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID            int32 `json:"id"`
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID   int32 `json:"to_account_id"`
	// must be positive
	Amount    int32     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
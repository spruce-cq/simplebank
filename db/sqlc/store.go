package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries with transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create a new Store struct with *sql.DB
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// begin transaction
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		// roback transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// commit transaction
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int32 `json:"from_account_id,omitempty"`
	ToAccountID   int32 `json:"to_account_id,omitempty"`
	Amount        int32 `json:"amount,omitempty"`
}

// TransferTxResult contains one transfer record and two entry
type TransferTxResult struct {
	Transfer    `json:"transfer,omitempty"`
	FromAccount Account `json:"from_account,omitempty"`
	ToAccount   Account `json:"to_account,omitempty"`
	FromEntry   Entry   `json:"from_entry,omitempty"`
	ToEntry     Entry   `json:"to_entry,omitempty"`
}

// TransferTx performs a money transfer from one to the other
// It creates a transfer record, and account entries, and update accounts's balance with a single db transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}); err != nil {
			return err
		}

		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}); err != nil {
			return err
		}
		// there are special dead lock when use sql as "select * frmo accounts where id = $id for update" to select account
		// then update account balance, when do this the select clouse will wait "create transfer" to release Exclusive lock to get shared lock
		// when other gorountine select for update aim to get share lock, this will be cause dead lock
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{Amount: -arg.Amount, ID: arg.FromAccountID})

		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{Amount: arg.Amount, ID: arg.ToAccountID})

		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store prodives all functions to execute db queries and transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute db queries with transaction
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore create a new Store struct with *sql.DB
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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

		// there are also a potential dead lock in the update order
		// like this: transaction1: "update account1 balance; update account2 balance", transaction2: "update account2 balance; update account1 balance"
		// so you should update in some order
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = s.addMoney(ctx, amountTip{arg.FromAccountID, -arg.Amount}, amountTip{arg.ToAccountID, arg.Amount})
		} else {
			result.ToAccount, result.FromAccount, err = s.addMoney(ctx, amountTip{arg.FromAccountID, -arg.Amount}, amountTip{arg.ToAccountID, arg.Amount})
		}

		return nil
	})
	return result, err
}

type amountTip struct {
	accountID int32
	amount    int32
}

func (s *SQLStore) addMoney(ctx context.Context, fromTip amountTip, toTip amountTip) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = s.AddAccountBalance(ctx, AddAccountBalanceParams{Amount: fromTip.amount, ID: fromTip.accountID})
	if err != nil {
		return
	}
	toAccount, err = s.AddAccountBalance(ctx, AddAccountBalanceParams{Amount: toTip.amount, ID: toTip.accountID})
	if err != nil {
		return
	}
	return
}

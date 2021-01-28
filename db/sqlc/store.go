package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore constructs a new Store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx safely executes a function with a database transaction.
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can not begin a transaction: %w", err)
	}

	// execute code
	q := New(tx)

	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("tx error: %v, rolback error: %w", err, errRb)
		}

		return fmt.Errorf("tx err: %w", err)
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx commit error: %w", err)
	}

	return nil
}

// TransferTxParams contains parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

// TransferTxResult contains result of the transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

// TransferTx performs a money transfer from the account to the another.
// It creates a new Transfer record with entries for both affected accounts and
// update their balances within a single database transaction.
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// transfer
		if result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg)); err != nil {
			return fmt.Errorf("failed to create a new transaction: %w", err)
		}

		// entries
		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}); err != nil {
			return fmt.Errorf("failed to create an entry for the sender: %w", err)
		}

		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}); err != nil {
			return fmt.Errorf("failed to create an entry for the receiver: %w", err)
		}

		// accounts
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: -arg.Amount,
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return fmt.Errorf("failed to modify (sender) account's balance: %w", err)
			}

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return fmt.Errorf("failed to modify (receiver) account's balance: %w", err)
			}
		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return fmt.Errorf("failed to modify (receiver) account's balance: %w", err)
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: -arg.Amount,
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return fmt.Errorf("failed to modify (sender) account's balance: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return TransferTxResult{}, fmt.Errorf("can not make a transaction: %w", err)
	}

	return result, nil
}

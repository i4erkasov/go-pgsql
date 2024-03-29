package pgx

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/i4erkasov/go-pgsql/pgxpool"
	"github.com/jackc/pgx/v4"
)

// Tx is an alias to pgx.Tx
type Tx = pgx.Tx

// Row is an alias to pgx.Row
type Row = pgx.Row

// Conn defines the interface for connection
type Conn interface {
	Begin(ctx context.Context) (Tx, error)
}

// TxManager defines the interface for transaction management.
type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context, Tx) error) error
	WithNestedTx(ctx context.Context, tFunc func(context.Context, Tx) error) error
}

// Transactor is a concrete implementation of TxManagerInterface using pgxpool.
type Transactor struct {
	conn Conn
}

// NewTxManager creates a new instance of TxManager with a given registry.
// It uses the master connection pool for managing transactions.
func NewTxManager(registry *pgxpool.Registry) (TxManager, error) {
	pools, err := registry.Pools()
	if err != nil {
		return nil, err
	}

	return &Transactor{
		pools.Master(),
	}, nil
}

type txContextKey struct{}
type txContextCounterKey struct{}

var (
	// txKey is a private key used for storing transaction context.
	txKey = txContextKey{}

	// txCounterKey is used for storing the transaction counter in the context.
	txCounterKey = txContextCounterKey{}

	// ErrNoTransaction is the error used when no transaction is found in the context.
	ErrNoTransaction = errors.New("no transaction in context")
)

// Begin starts a new transaction and stores it in the context.
func (t *Transactor) begin(ctx context.Context) (context.Context, error) {
	tx, err := t.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, txKey, tx), nil
}

// Commit commits the transaction stored in the context.
func (t *Transactor) commit(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(Tx)
	if !ok {
		return ErrNoTransaction
	}
	return tx.Commit(ctx)
}

// Rollback aborts the transaction stored in the context.
func (t *Transactor) rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(Tx)
	if !ok {
		return ErrNoTransaction
	}
	return tx.Rollback(ctx)
}

// WithTx executes a function within the context of a transaction.
// This method checks if there is already an ongoing transaction in the context.
// If not, it starts a new transaction and then executes the provided function.
// After the function execution, it commits the transaction if no errors occurred,
// or rollbacks in case of an error or panic.
// The transaction object is passed to the function, allowing direct transaction control.
func (t *Transactor) WithTx(ctx context.Context, tFunc func(context.Context, Tx) error) (err error) {
	// Check if there is already a transaction in the context
	tx, ok := ctx.Value(txKey).(Tx)
	if !ok {
		// Start a new transaction if there isn't one
		ctx, err = t.begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			// Handle the end of the transaction
			if p := recover(); p != nil {
				_ = t.rollback(ctx)
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				_ = t.rollback(ctx) // err is non-nil; rollback the transaction
			} else {
				err = t.commit(ctx) // err is nil; if Commit returns error update err
			}
		}()
		// Get the transaction object after beginning a new transaction
		tx = ctx.Value(txKey).(Tx)
	}

	// Execute the function passed, providing the transaction object
	err = tFunc(ctx, tx)

	return err
}

// WithNestedTx executes a function within the context of a potentially nested transaction.
// It manages transaction nesting using a counter to track the depth of nested transactions.
// If there is no active transaction, it starts a new one.
func (t *Transactor) WithNestedTx(ctx context.Context, tFunc func(context.Context, Tx) error) (err error) {
	// Start a new transaction or increment the nested transaction counter.
	ctx, nested, err := t.beginNestedTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		// Decrement the transaction counter when exiting the function.
		ctx = tickTxCounter(ctx, -1)

		// Handle any panics or errors, and rollback if necessary.
		if p := recover(); p != nil {
			if nested {
				_ = t.rollback(ctx)
			}
			panic(p)
		} else if err != nil {
			if nested {
				_ = t.rollback(ctx)
			}
		} else {
			if nested {
				err = t.commit(ctx)
			}
		}
	}()

	// Get the transaction object from the context to pass to the function.
	tx := ctx.Value(txKey).(Tx)

	// Execute the provided function within the transaction context.
	err = tFunc(ctx, tx)

	return err
}

// beginNestedTx starts a new transaction if one is not already in progress,
// and manages the transaction counter. If an existing transaction is detected,
// it increments the counter without starting a new transaction.
func (t *Transactor) beginNestedTx(ctx context.Context) (context.Context, bool, error) {
	// Increment the transaction counter and check if a transaction is already in progress.
	ctx = tickTxCounter(ctx, 1)
	if _, ok := ctx.Value(txKey).(Tx); ok {
		// If a transaction is already in progress, use it without starting a new one.
		return ctx, false, nil
	}

	// Start a new transaction if there isn't one already.
	var err error
	ctx, err = t.begin(ctx)
	if err != nil {
		// In case of an error, decrement the counter back.
		ctx = tickTxCounter(ctx, -1)
		return ctx, false, err
	}
	return ctx, true, nil
}

// tickTxCounter safely increments or decrements the transaction counter in the context.
// Returns the updated context.
func tickTxCounter(ctx context.Context, val int32) context.Context {
	count, ok := ctx.Value(txCounterKey).(*int32)

	if !ok {
		// Initialize the counter if it's not present in the context.
		var cnt int32
		count = &cnt
		ctx = context.WithValue(ctx, txCounterKey, count)
	}

	// Atomically update the counter.
	atomic.AddInt32(count, val)

	return ctx
}

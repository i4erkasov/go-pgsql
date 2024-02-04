package pgx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type txManagerTestSuite struct {
	suite.Suite
}

func (t *txManagerTestSuite) TestWithTxSuccess() {
	t.T().Parallel()

	connMock := new(ConnMock)
	txMock := new(TxMock)

	// Setup: Expect Begin to be called, returning mock transaction and nil as an error.
	connMock.On("Begin", mock.Anything).Return(txMock, nil) // Использование mock.Anything вместо mock.AnythingOfType("*context.emptyCtx")
	// Expect Commit to be called due to successful function execution.
	txMock.On("Commit", mock.Anything).Return(nil) // Использование mock.Anything

	transactor := Transactor{conn: connMock}

	err := transactor.WithTx(context.Background(), func(ctx context.Context, tx Tx) error {
		// Your business logic here
		return nil
	})

	// Assert no error was returned and expectations were met.
	assert.NoError(t.T(), err)
	connMock.AssertExpectations(t.T())
	txMock.AssertExpectations(t.T())
}

func (t *txManagerTestSuite) TestWithTxInsertUser() {
	t.T().Parallel()

	connMock := new(ConnMock)
	txMock := new(TxMock)

	// Setup: expect the Begin method to be called, returning a mock transaction and nil as an error.
	connMock.On("Begin", mock.Anything).Return(txMock, nil)
	// Expect the Exec method to be called with an SQL insert query, returning nil as an error.
	// Ensure that the arguments are passed correctly.
	txMock.On("Exec", mock.Anything, "INSERT INTO users (name, email) VALUES ($1, $2)", "John Doe", "john@example.com").Return(nil, nil) // Note the addition of nil for the pgconn.CommandTag
	// Expect the transaction to be committed.
	txMock.On("Commit", mock.Anything).Return(nil)

	transactor := Transactor{conn: connMock}

	err := transactor.WithTx(context.Background(), func(ctx context.Context, tx Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John Doe", "john@example.com")
		return err
	})

	// Verify that no errors occurred and all expectations were met.
	assert.NoError(t.T(), err)
	connMock.AssertExpectations(t.T())
	txMock.AssertExpectations(t.T())
}

func (t *txManagerTestSuite) TestWithTxRollbackOnError() {
	t.T().Parallel()

	connMock := new(ConnMock)
	txMock := new(TxMock)

	// Setup: Expect Begin to be called, returning mock transaction and nil as an error.
	connMock.On("Begin", mock.Anything).Return(txMock, nil) // Использование mock.Anything
	// Expect Rollback to be called due to an error in the function.
	txMock.On("Rollback", mock.Anything).Return(nil) // Использование mock.Anything

	transactor := Transactor{conn: connMock}

	err := transactor.WithTx(context.Background(), func(ctx context.Context, tx Tx) error {
		// Simulating an error in the business logic
		return errors.New("error in transaction")
	})

	// Assert an error was returned and expectations were met.
	assert.Error(t.T(), err)
	connMock.AssertExpectations(t.T())
	txMock.AssertExpectations(t.T())
}

// TestWithNestedTxSuccess tests successful execution within a nested transaction context.
func (t *txManagerTestSuite) TestWithNestedTxSuccess() {
	t.T().Parallel()

	connMock := new(ConnMock)
	txMock := new(TxMock)

	// Expect Begin to be called, returning mock transaction and nil as an error.
	connMock.On("Begin", mock.Anything).Return(txMock, nil)
	// Expect Commit to be called due to successful function execution within the nested transaction.
	txMock.On("Commit", mock.Anything).Return(nil)

	transactor := Transactor{conn: connMock}

	err := transactor.WithNestedTx(context.Background(), func(ctx context.Context, tx Tx) error {
		// Simulate successful business logic within the nested transaction.
		return nil
	})

	// Assert no error was returned and expectations for both Begin and Commit were met.
	assert.NoError(t.T(), err)
	connMock.AssertExpectations(t.T())
	txMock.AssertExpectations(t.T())
}

// TestWithNestedTxRollbackOnError tests that a transaction is correctly rolled back on error within the nested transaction context.
func (t *txManagerTestSuite) TestWithNestedTxRollbackOnError() {
	t.T().Parallel()

	connMock := new(ConnMock)
	txMock := new(TxMock)

	// Expect Begin to be called, returning mock transaction and nil as an error.
	connMock.On("Begin", mock.Anything).Return(txMock, nil)
	// Expect Rollback to be called due to an error in the function executed within the nested transaction.
	txMock.On("Rollback", mock.Anything).Return(nil)

	transactor := Transactor{conn: connMock}

	err := transactor.WithNestedTx(context.Background(), func(ctx context.Context, tx Tx) error {
		// Simulate an error in the business logic within the nested transaction.
		return errors.New("error in nested transaction")
	})

	// Assert an error was returned and expectations for both Begin and Rollback were met.
	assert.Error(t.T(), err)
	connMock.AssertExpectations(t.T())
	txMock.AssertExpectations(t.T())
}

func (t *txManagerTestSuite) TestWithNestedTxConditionalLogic() {
	t.T().Parallel()

	connMock := new(ConnMock)
	txMock := new(TxMock)
	rowMock := new(RowMock)

	// Configure the behavior of the Scan method to simulate successful retrieval of userID.
	// Important: configure rowMock to correctly populate the given variable, for example, userID = 1 to simulate an existing user.
	rowMock.On("Scan", mock.AnythingOfType("*int")).Return(func(dest ...interface{}) error {
		*dest[0].(*int) = 1 // Simulate the presence of a user with ID = 1
		return nil
	}).Once()

	connMock.On("Begin", mock.Anything).Return(txMock, nil)
	// Correctly set the return value for QueryRow: return rowMock to simulate the query result.
	txMock.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(rowMock)
	txMock.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil, nil)
	txMock.On("Commit", mock.Anything).Return(nil)

	transactor := Transactor{conn: connMock}

	err := transactor.WithNestedTx(context.Background(), func(ctx context.Context, tx Tx) error {
		var userID int
		if err := tx.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", "john@example.com").Scan(&userID); err != nil {
			// Assume if there's an error, the user doesn't exist and needs to be added.
			_, err = tx.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John Doe", "john@example.com")
			return err
		} else {
			// If the user exists, update their details.
			_, err = tx.Exec(ctx, "UPDATE users SET name = $1 WHERE email = $2", "John Updated", "john@example.com")
			return err
		}
	})

	// Verify that no errors occurred and all expectations were met.
	assert.NoError(t.T(), err)
	connMock.AssertExpectations(t.T())
	txMock.AssertExpectations(t.T())
	rowMock.AssertExpectations(t.T())
}

func TestTxManager_Run(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(txManagerTestSuite))
}

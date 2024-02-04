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

func TestTxManager_Run(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(txManagerTestSuite))
}


# Transaction Management with Transactor

`Transactor` provides robust transaction management capabilities, allowing you to execute operations within transactions using `WithTx` and manage nested transactions with `WithNestedTx`. Below are the use cases and examples for each method.

## Using `WithTx`

The `WithTx` method is used for executing operations within a single transaction. It's ideal for operations that need to be atomic to maintain data integrity.

### Use Cases

- **Single Database Action**: Use `WithTx` for individual insert, update, or delete operations.
- **Read and Update**: Useful for scenarios where you read data and then perform updates based on that data.
- **Aggregate Operations**: For executing multiple operations that need to be atomic.

### Example

```go
func UpdateUserEmail(ctx context.Context, txManager pgx.TxManager, userID int, newEmail string) error {
    return txManager.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
        // Perform the update operation
        _, err := tx.Exec(ctx, "UPDATE users SET email = $1 WHERE id = $2", newEmail, userID)
        return err
    })
}
```

## Using `WithNestedTx`

`WithNestedTx` is designed for managing potentially nested transactions. It's useful when operations are complex, involve conditional logic, or when you want to maintain modularity in your transaction logic.

### Use Cases

- **Nested Operations with Conditions**: When different operations need to be performed based on conditions, within their own transactional scope.
- **Recursive Operations**: For operations that call themselves with different parameters, where each call should be in its own transaction.
- **Modular Business Processes**: When a business process is divided into modules that can execute independently but within one overarching transaction.
- **Long Transactional Chains**: For a series of logically connected operations that should be rolled back together in case of an error.

### Example

```go
func ProcessOrder(ctx context.Context, txManager pgx.TxManager, orderID int) error {
    return txManager.WithNestedTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
        // Step 1: Update order status
        if err := updateOrderStatus(ctx, tx, orderID, "processing"); err != nil {
            return err
        }

        // Step 2: Reserve inventory
        if err := reserveInventory(ctx, tx, orderID); err != nil {
            return err
        }

        // Step 3: Process payment
        return processPayment(ctx, tx, orderID)
    })
}
```

These methods offer flexibility and control over transaction management, ensuring data integrity and consistency across your application. Use `WithTx` for straightforward transactional operations and `WithNestedTx` for more complex or conditional transaction logic.

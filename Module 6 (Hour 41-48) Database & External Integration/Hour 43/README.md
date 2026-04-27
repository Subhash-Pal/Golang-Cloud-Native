# Hour 43 - Transactions and Rollback

This project shows how database transactions protect data consistency.

## Flow

1. Seed two accounts
2. Start a transaction
3. Debit one account
4. Simulate a failure so the transaction rolls back
5. Run the transfer again successfully

## Run

```powershell
go mod tidy
.\run.ps1
```

## Table Used

The example creates and uses the `hour43_accounts` table.

## Requirement

PostgreSQL must be running before you start this example.

If you use Docker for PostgreSQL, start Docker Desktop first, then start the container before running `.\run.ps1`.

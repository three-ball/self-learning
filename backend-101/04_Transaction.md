# 04. Transaction & Locking

> A transaction is a unit of work that you want to treat as "a whole." It has to either happen in full or not at all.

## Introduction

- A transaction is a group of SQL queries that are treated as atomically, as a single unit of work.
- All or nothing.

```sql
BEGIN;
SELECT * FROM account_balance WHERE account_id = 1;

UPDATE account_balance SET balance = balance - 100 WHERE account_id = 1;
UPDATE account_balance SET balance = balance + 100 WHERE account_id = 2;

SELECT * FROM account_balance WHERE account_id IN (1, 2);
-- If everything goes well, you can commit the transaction
COMMIT;

-- If something goes wrong, you can roll back the transaction
ROLLBACK;
```

### 1.1. ACID

- **Atomicity**: All operations in a transaction are treated as a single unit. If one operation fails, the entire transaction fails.
- **Consistency**: Data should be valid according to all predefined rules, including constraints, cascades, and triggers.
- **Isolation**: Transactions do not affect each other. Each transaction is isolated from others until it is committed.
- **Durability**: Once a transaction is committed, it remains so, even in the event of a system failure. The changes made by the transaction are permanent.

### 1.2. Isolation Levels

#### 1.2.1. Read Uncommitted

- Transactions can read data that is not yet committed by other transactions -> Dirty Reads.
- Example: If Transaction A is updating a row, Transaction B can read the uncommitted changes.
- Rarely used due to the risk of dirty reads.

#### 1.2.2. Read Committed (Default of PostgreSQL)

- Transactions can only read data that has been committed by other transactions.
- Prevents dirty reads.
- Problem: Non-repeatable reads can occur. If Transaction A reads a row, and Transaction B updates it before Transaction A reads it again, Transaction A will see different data.

#### 1.2.3. Repeatable Read (Default of MySQL)

- Transactions can read the same data multiple times and get the same result.
- Prevents dirty reads and non-repeatable reads.
- Problem: Phantom reads can occur. Some records "appear" and "disappear" between reads in the same transaction. Example: If Transaction A reads a set of rows, and Transaction B inserts a new row that matches the criteria before Transaction A reads again, Transaction A will see different results.

#### 1.2.4. Serializable

- The highest isolation level.
- `Serializable` transactions are executed in such a way that they appear to be executed one after the other, even if they are actually executed concurrently.
- In a nutshell, it prevents dirty reads, non-repeatable reads, and phantom reads. `Serializable` places a lock on every row it reads.
- Problem: It can lead to performance issues due to increased locking and blocking, especially in high-concurrency environments.

## 2. How Transactions Work

### 2.1. Transaction Logging

- Write Operation: WAL (Write-Ahead Logging) is a common technique used by storage engines to ensure durability and consistency of transactions.
    - Write-ahead logging: The storage engine can then write the change to the transaction log on disk before applying it to the database.
    - The storage engine can change the data in memory and then write it to the transaction log.
    - Later, the storage engine can apply the changes to the database.

### 2.2. Multi-Version Concurrency Control (MVCC)

- InnoDB uses MVCC and undo logs to accomplish A, C, I in ACID.
- Multiversion concurrency control means that changes to a row create a new version of the row.
- When the original transaction commits and no other active transactions are using the old version, the old version is deleted.
- Undo logs are saved in the InnoDB buffer pool.
- A transaction can hold many locks and undo logs, so it will impact the performance of the database.

### 2.3. Some problems:

- **Large Transactions**: modify too many rows, can lead to performance issues due to the amount of data being processed and the number of locks held.
- **Long Transactions**: Long running transaction, can lead to contention and blocking, as other transactions may be waiting for locks held by the long transaction.
- **Abandoned Transactions**: Client connection vanished during actiove transaction. If a transaction is not committed or rolled back, it can lead to locks being held indefinitely, causing contention and blocking for other transactions.

**Notes**:
- `Transaction` is different from `locking`. `Transactions` are a way to group multiple operations into a single unit of work, while `locking` is a mechanism to control access to resources in a concurrent environment.

## 3. Locking

### 3.1. Problem Context

- Context: Fund transfer between two accounts.
- Problem: If two transactions try to transfer money between the same two accounts at the same time, they can interfere with each other and lead to incorrect results.

### 3.2. Types of Locks (in MySQL)

**Classification of locks in MySQL:**
- **Shared Lock (S Lock)**: Allows a transaction to read a row but not modify it. Other transactions can also acquire shared locks on the same row.
- **Exclusive Lock (X Lock)**: Lock to read and write a row.
- Example 1:
    - Transaction A locks row X for reading (S Lock).
    - Transaction B can also lock row X for reading (S Lock).
- Example 2:
    - Transaction A locks row X for reading (S Lock).
    - Transaction B tries to lock row X for writing (X Lock) but is blocked until Transaction A releases the lock.
- Example 3:
    - Transaction A locks row X for writing (X Lock).
    - Transaction B can't lock row X for reading or writing until Transaction A releases the lock.
- Solve problem above:

```sql
-- First transaction
begin;
select * from account_balance where account_id = 1 for update;
update account_balance set balance = balance - 100 where account_id = 1;
commit;
-- Second transaction
begin;
select * from account_balance where account_id = 1 for update;
update account_balance set balance = balance - 200 where account_id = 1;
commit;
```

- **Row Lock Types:**

| Lock | Abbreviation | Locks Gap | Description |
|------|--------------|-----------|-------------|
| Record Lock | REC_NOT_GAP | No | Locks a specific row in a table. |
| Gap Lock | GAP | Yes | Locks the gap before (less than) a specific row in a table. |
| Next-Key Lock | NEXT | Yes | Locks the gap before a specific row and the row itself. |
| Insert Intention Lock | INSERT_INTENTION |  | Allows INSERT into a gap without blocking other transactions. |
 

### 3.3. Lock Contention

- **Lock Contention**: Occurs when multiple transactions try to acquire locks on the same resource, leading to delays and potential deadlocks.

### 3.4. Pessimistic vs Optimistic Locking (classificafy by implementation)

- **Pessimistic Locking**: Locks resources as soon as they are accessed, preventing other transactions from modifying them until the lock is released.
    - Example: Using `SELECT ... FOR UPDATE` to lock rows for update.
    - **Pros:** Prevents conflicts and ensures data integrity.
    - **Cons:** Can lead to contention and reduced concurrency. Can occur bottlenecks if many transactions are waiting for locks.
- **Optimistic Locking**: Assumes that conflicts are rare and allows transactions to proceed without locking resources initially. It checks for conflicts before committing. Implement using Compare-And-Swap (CAS) or versioning.
    - Example: Using a version column to check if the row has been modified before committing.
    - **Pros:** Higher concurrency and less contention.
    - **Cons:** Requires additional checks before committing, which can lead to retries if conflicts are detected.

### 3.5. Best Practices

**Subtract first, then add.**
    **Why this matters**:
    - If account 1 has $50 and we try to transfer $100, subtracting first will fail immediately
    - Adding first might temporarily create money that shouldn't exist
    - Helps with constraint validation (e.g., balance >= 0)

**Use record lock by leveraging primary key:** Primary key lookups are faster and create more precise locks compared to secondary indexes.
    **Why primary keys are better**:
    - Direct access via clustered index
    - Guaranteed unique, so only one row locked
    - Faster lookup and lock acquisition

```sql
-- ✅ CORRECT: Use primary key for locking
BEGIN;
SELECT * FROM users WHERE id = 123 FOR UPDATE;  -- Fast, precise record lock
UPDATE users SET balance = balance - 100 WHERE id = 123;
COMMIT;

-- ❌ SUBOPTIMAL: Use secondary index for locking
BEGIN;
SELECT * FROM users WHERE email = 'john@example.com' FOR UPDATE;  -- Slower, may lock more rows
UPDATE users SET balance = balance - 100 WHERE email = 'john@example.com';
COMMIT;
```

**Avoid using secondary index for locking.**

**Example**:
```sql
-- ✅ CORRECT: Lock by primary key
BEGIN;
SELECT * FROM orders WHERE order_id = 12345 FOR UPDATE;
UPDATE orders SET status = 'processed' WHERE order_id = 12345;
COMMIT;

-- ❌ PROBLEMATIC: Lock by secondary index
BEGIN;
SELECT * FROM orders WHERE customer_id = 678 FOR UPDATE;  -- May lock multiple rows
UPDATE orders SET status = 'processed' WHERE customer_id = 678;
COMMIT;
```

**Examine row locks.**
    **What to monitor**:
    - Lock wait times
    - Deadlock frequency
    - Lock contention patterns
    - Long-running transactions

**Leverage Optimistic Locking**: Use version numbers or timestamps to detect concurrent modifications without holding locks.
- Combine pessimistic and optimistic locking.
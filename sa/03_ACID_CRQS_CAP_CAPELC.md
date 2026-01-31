# ACID - CQRS - CAP - PCAELC

## ACID

```mermaid
mindmap
  root((ACID))
    Atomicity
      All-or-nothing
      Either all succeed
      Or none applied
      Rollback on failure
    Consistency
      Valid state before & after
      Adhere to rules
      Constraints & triggers
    Isolation
      Read Phenomena
        Dirty Read
        Non-repeatable Read
        Phantom Read
      Isolation Levels
        Read Uncommitted
        Read Committed
        Repeatable Read
        Serializable
    Durability
      Permanent changes
      Survives system failure
      Once committed
```

* ACID stands for `Atomicity`, `Consistency`, `Isolation`, `Durability`.

### Atomicity

* `Atomicity` means a transaction is **all-or-nothing** either all its operations succeed, or none are applied. If any part fails, the entire transaction is rolled back to keep the database consistent.

### Consistency

* `Consistency` in transactions means that the database **must remain in a valid state before and after a transaction**. Any data written must adhere to all defined rules, constraints, and triggers.

### Isolation

* `Isolation` ensures that **concurrent transactions do not interfere with each other**. Each transaction should operate as if it is the only one running, preventing issues like dirty reads or lost updates.
* Two main things to consider: Read phenomena and Isolation levels.
* Read phenomena:
  * **Dirty Read**: Reading uncommitted data from another transaction.
  * **Non-repeatable Read**: Reading the same data multiple times within a transaction and getting different results because another transaction modified it.
  * **Phantom Read**: Reading a set of rows that satisfy a condition, but another transaction inserts or deletes rows that affect the result set.
* Isolation levels:
  * **Read Uncommitted**: Lowest level but highest concurrency. Allows dirty reads: Reads data that has been modified but not yet committed by other transactions.
  * **Read Committed**: Prevents dirty reads. A transaction **can only read data that has been committed** by other transactions.
  * **Repeatable Read**: This isolation level guarantees that a transaction will see the same data throughout its duration, even if other transactions commit changes to the data. As it only sees committed data before the transaction started, it prevents dirty reads and non-repeatable reads. However, new rows may still be inserted that satisfy the query condition, which can lead to phantom reads.
  * **Serializable**: Highest isolation level. Transactions are completely isolated from each other, **as if they were executed serially**. Prevents dirty reads, non-repeatable reads, and phantom reads.

| Isolation Level  | Dirty Read | Non-repeatable Read | Phantom Read | Desc             |
| ---------------- | ---------- | ------------------- | ------------ | -------------------------- |
| READ UNCOMMITTED | YES       | YES                | YES         | Read uncommitted data      |
| READ COMMITTED   | NO    | YES                | YES         | Read committed data        |
| REPEATABLE READ  | NO    | NO             | YES         | Repeatable read data       |
| SERIALIZABLE     | NO    | NO             | NO      | Results as if executed serially |

### Durability

* `Durability` ensures that once a transaction is committed, **its changes are permanently saved, even if the system fails**.

## CQRS

```mermaid
flowchart LR
    Client[Client / API Gateway]

    Client -->|Commands| WriteDB[(Write Database)]
    Client -->|Queries| ReadDB[(Read Database)]
    WriteDB -->|Eventual Consistency| ReadDB
```

* CQRS stands for `Command Query Responsibility Segregation`.
* It is a design pattern that separates the read and write operations of a data store into different models.
* **CQRS often goes hand-in-hand with Event Sourcing**, where state changes are stored as a sequence of events. [Event sourcing, CQRS, stream processing and Apache Kafka: What’s the connection?](https://www.confluent.io/blog/event-sourcing-cqrs-stream-processing-apache-kafka-whats-connection/)

![CQRS Diagram](images/03_confluent_cqrs.jpeg)

## CAP Theorem

![CAP Theorem](images/03_cap_theorem.png)

* CAP stands for `Consistency`, `Availability`, `Partition Tolerance`.
* CAP theorem states that in a distributed data store, you can only achieve two out of the three guarantees at the same time.
* CA only works in non-distributed systems (when all data is stored on a single node), but in practice, distributed systems must handle network partitions, so CP and AP are the more relevant trade-offs.
* Banking, financial systems often prioritize Consistency and Partition Tolerance (CP) to ensure accurate transactions, while social media platforms may prioritize Availability and Partition Tolerance (AP) to ensure users can always access content, even if some data is temporarily inconsistent.

## PCAELC

```mermaid
mindmap
  root((PCAELC))
    Partition_Tolerance
      Yes
        Consistency
        Availability
      No
        Latency
        Consistency
```

* PACELC theorem is an extension of the CAP theorem that adds considerations for latency during normal operation (when there are no network partitions).
* It states that in the event of a network partition (P), a distributed system must choose between Consistency (C) and Availability (A). Else (E), when the system is running normally, it must choose between Latency (L) and Consistency (C).
* The trade-offs between "Latency" and "Consistency" during normal operation depend on the specific requirements of the application. When we need low latency, we might accept some level of eventual consistency. Conversely, if strong consistency is crucial, we might tolerate higher latency due to synchronization overhead.
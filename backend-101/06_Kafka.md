# Kafka

> Kafka is a distributed event streamin platform. Kafka was originally developed by LinkedIn and later open-sourced in 2011. It is now part of the Apache Software Foundation.

## Introduction

```mermaid
graph TD
    %% Publisher
    P[📤 Order Service<br/>Producer] 
    
    %% Kafka Cluster
    subgraph KC[" Kafka Cluster "]
        subgraph T[" Topic: order-events "]
            P0[Partition 0<br/>📨 📨 📨]
            P1[Partition 1<br/>📨 📨 📨]
            P2[Partition 2<br/>📨 📨 📨]
        end
    end
    
    %% Consumers
    C1[📥 Inventory Service<br/>Consumer 1]
    C2[📥 Payment Service<br/>Consumer 2]
    C3[📥 Notification Service<br/>Consumer 3]
    C4[📥 Analytics Service<br/>Consumer 4]
    
    %% Connections
    P -->|Publish Messages| T
    T -->|Subscribe & Consume| C1
    T -->|Subscribe & Consume| C2
    T -->|Subscribe & Consume| C3
    T -->|Subscribe & Consume| C4
    
    %% Dark Theme Material Design Styling
    style P fill:#1976D2,stroke:#BBDEFB,stroke-width:2px,color:#FFFFFF
    style KC fill:#212121,stroke:#757575,stroke-width:2px,color:#FFFFFF
    style T fill:#424242,stroke:#9E9E9E,stroke-width:2px,color:#FFFFFF
    style P0 fill:#37474F,stroke:#90A4AE,stroke-width:1px,color:#FFFFFF
    style P1 fill:#37474F,stroke:#90A4AE,stroke-width:1px,color:#FFFFFF
    style P2 fill:#37474F,stroke:#90A4AE,stroke-width:1px,color:#FFFFFF
    style C1 fill:#388E3C,stroke:#A5D6A7,stroke-width:2px,color:#FFFFFF
    style C2 fill:#F57C00,stroke:#FFD54F,stroke-width:2px,color:#FFFFFF
    style C3 fill:#7B1FA2,stroke:#E1BEE7,stroke-width:2px,color:#FFFFFF
    style C4 fill:#D32F2F,stroke:#FFCDD2,stroke-width:2px,color:#FFFFFF
    
    %% Link styling for dark theme
    linkStyle 0 stroke:#64B5F6,stroke-width:3px
    linkStyle 1 stroke:#81C784,stroke-width:3px
    linkStyle 2 stroke:#FFB74D,stroke-width:3px
    linkStyle 3 stroke:#CE93D8,stroke-width:3px
    linkStyle 4 stroke:#EF5350,stroke-width:3px
```

- Distributed streaming platform. high throughput with long-term retention.
- Pub/sub pattern, fan-out mechanism.
- Used for building real-time data pipelines and streaming applications.
- Asynchronous communication between services (different from synchronous HTTP requests, where client waits for a response).
- Kafka can be used to decouple services, allowing them to operate independently. This lead to better `scalability` and `fault tolerance` (if one service fails, others can continue to operate).
- **Pros**:
  - High throughput.
  - Asynchronous processing.
  - Decoupled architecture.
  - Scalable and fault-tolerant.
- **Cons**:
  - Complexity in setup and management.
  - Operational overhead.
  - Cost.

### The difference between `Event` vs `Request/Response`

- **Event**: Just a thing that happened. It may a bussiness fact that value to more than one service. Event does not require a response.
- **Request/Response**: A request is made to a service, and a response is expected. It is a synchronous communication pattern.

### Event Streaming

- **Streaming**: Continuous flow of data (or events) that can be processed in real-time.
- **Event streaming** includes: How to get events in/out, how to store events, how to order them, etc.

> If Event streaming is good, why not use it everywhere?

- Event streaming is not always the best solution for every use case. **It always has its own trade-offs**:
  - Cost and complexity.
  - Error handling and debugging can be more complex.
  - Operational overhead.
  - Learning curve for developers.
  - Some usecase expect a response, which is not suitable for event streaming!

## Core Concepts

- **Broker**: A Kafka server that stores and serves messages. One cluster can have multiple brokers.
- **Zookeeper**: Mangages broker and mananges the overall cluster state. It is used for leader election, configuration management, and coordination.
- **Kafka Raft (KRaft)**: An alternative to Zookeeper for managing Kafka clusters. It is a newer approach that simplifies the architecture by removing the dependency on Zookeeper.

```mermaid
graph TD
    subgraph ZK[" Zookeeper Cluster "]
        Z1[🗂️ Zookeeper 1]
        Z2[🗂️ Zookeeper 2] 
        Z3[🗂️ Zookeeper 3]
    end
    
    subgraph KC[" Kafka Cluster "]
        B1[🏢 Broker 1]
        B2[🏢 Broker 2]
        B3[🏢 Broker 3]
    end
    
    %% Zookeeper manages brokers
    ZK -.->|Coordinates| KC
    
    style ZK fill:#FFE0B2,stroke:#FF8F00
    style KC fill:#E8F5E8,stroke:#4CAF50
```

- **Topic**: A category or feed name to which messages are published. Topics are partitioned for scalability.
- **Partition**: A topic can be divided into multiple partitions, which allows for parallel processing and scalability. Each partition is an ordered, immutable sequence of messages. kafka uses partition to scale horizontally. Partitions act as a queue, where messages are stored in the order they are received.

```mermaid
graph TD
    subgraph T1[" Topic: user-events "]
        P1[📂 Partition 0<br/>Msg1 → Msg2 → Msg3]
        P2[📂 Partition 1<br/>Msg4 → Msg5 → Msg6]
        P3[📂 Partition 2<br/>Msg7 → Msg8 → Msg9]
    end
    
    %% Producers write to partitions
    PROD[📤 Producer] -->|Hash Key| P1
    PROD -->|Hash Key| P2
    PROD -->|Hash Key| P3
    
    %% Consumers read from partitions
    P1 --> CONS1[📥 Consumer 1]
    P2 --> CONS2[📥 Consumer 2]
    P3 --> CONS3[📥 Consumer 3]
    
    style T1 fill:#E8F5E8,stroke:#4CAF50,color:#000000
    style PROD fill:#1976D2,stroke:#BBDEFB,color:#FFFFFF
    style CONS1 fill:#388E3C,stroke:#A5D6A7,color:#FFFFFF
    style CONS2 fill:#F57C00,stroke:#FFD54F,color:#FFFFFF
    style CONS3 fill:#7B1FA2,stroke:#E1BEE7,color:#FFFFFF
```

- **Offset** is a unique identifier for each message within a partition. It is an integer that represents the position of the message in the partition. Offsets are used to track which messages have been consumed by a consumer.
- **Record**: A message in Kafka is called a record. Record has 6 attributes:
  - `Key`: Optional identifier for the record, used for partitioning or grouping related records.
  - `Value`: The actual data of the record. Nullable.
  - `Headers`: Optional metadata associated with the record.
  - `Partition`: The partition to which the record belongs.
  - `Offset`: The unique identifier for the record within the partition.
  - `Timestamp`: The time when the record was produced. In milliseconds since epoch.
- **Producer**: An application that publishes messages to Kafka topics. Producers can send messages to specific partitions based on the key or round-robin distribution.
- **Consumer**: An application that subscribes to Kafka topics and processes messages. Consumers can read messages from one or more partitions.

```mermaid
graph LR
    P[📤 Producer] -->|Publish| T[📋 Topic A] -->|Consume| C[📥 Consumer]
    
    style P fill:#1976D2,stroke:#BBDEFB,color:#FFFFFF
    style T fill:#4CAF50,stroke:#A5D6A7,color:#FFFFFF
    style C fill:#FF9800,stroke:#FFD54F,color:#FFFFFF
```

- **Consumer Group**: A group of consumers that work together to consume messages from a topic. Each consumer in the group reads from a different partition, allowing for parallel processing. Each message is consumed by only one consumer in the group.
  - The first consumer in a group will read all partitions.
  - When the second consumer joins the group, **a rebalance is triggered**, and partitions are reassigned to consumers in the group. This action is automatic and transparent to the application.
  - A partition is assigned to only one consumer in a group at a time. There may be more consumers than partitions, in which case some consumers will not receive any messages. The redundant consumers will be idle and can be used for failover.
  - Different consumer groups do not interfere with each other. Each group can consume the same messages independently.

### Consumer Groups

```mermaid
graph TD
    subgraph T[" Topic: order-events "]
        P0[📂 Partition 0<br/>Order1 → Order2 → Order3]
        P1[📂 Partition 1<br/>Order4 → Order5 → Order6]
        P2[📂 Partition 2<br/>Order7 → Order8 → Order9]
    end
    
    subgraph CG1[" Consumer Group: processing-service "]
        C1[📥 Consumer 1]
        C2[📥 Consumer 2]
        C3[📥 Consumer 3]
    end
    
    subgraph CG2[" Consumer Group: analytics-service "]
        C4[📥 Consumer A]
        C5[📥 Consumer B]
    end
    
    %% Consumer Group 1 - Each consumer gets one partition
    P0 --> C1
    P1 --> C2
    P2 --> C3
    
    %% Consumer Group 2 - Fewer consumers than partitions
    P0 --> C4
    P1 --> C4
    P2 --> C5
    
    style T fill:#E8F5E8,stroke:#4CAF50,color:#000000
    style CG1 fill:#E3F2FD,stroke:#2196F3,color:#000000
    style CG2 fill:#FFF3E0,stroke:#FF9800,color:#000000
    style C1 fill:#1976D2,stroke:#BBDEFB,color:#FFFFFF
    style C2 fill:#1976D2,stroke:#BBDEFB,color:#FFFFFF
    style C3 fill:#1976D2,stroke:#BBDEFB,color:#FFFFFF
    style C4 fill:#F57C00,stroke:#FFD54F,color:#FFFFFF
    style C5 fill:#F57C00,stroke:#FFD54F,color:#FFFFFF
```

- **Partitioner**: A component that determines which partition a message should be sent:
  - If partition is specified, the message is sent to that partition.
  - If no partition is specified, the partitioner uses the key to determine the partition:
    - **If key is provided**, it hashes the key to determine the partition: `partion` = `murmur2(key)` % `(number_of_partitions - 1)`. Same key will always go to the same partition.
    - **If no key is provided**:
      - In kafka 2.4.0 earlier, it uses round-robin to distribute messages across partitions.
      - In kafka 2.4.0 and later, `Sticky Partitioner` is used. `Sticky Partitioner` improves the performance of the producer especially with high throughput. The producer sticky partitioner will:
        - Send messages to the same partition until the batch is full or the linger time is reached.
        - After that, it will switch to another partition.
        - This approach reduces the number of requests and improves throughput.
- **Replication**: A replica is a copy of a partition stored on another broker. Replication provides fault tolerance and high availability. Replicas are distributed across brokers to ensure that if one broker fails, the data is still available on other brokers.
  - The number of replicas is defined by the `replication factor` of the topic.
  - The default replication factor is 1, which means no replication. Replicattion factor should be greater than 1 and less than the size of the cluster.
- **Leader and Follower**:
  - Each partition has one leader and zero or more followers.
  - Producer only sends messages to the leader.
  - Consumers read from the leader replica by default.

## Producer

### ACK

- ACK is a signal sent from receiver indicating tht a mesasge has been received.
- ACKS is the number of acknowledgments the producer requires the leader to receive before considering a request complete.

### ACKS = 0

- The producer does not wait for any acknowledgment from the broker.
- Just fire and forget.
- **Pros**:
  - Fastest option, no waiting for acknowledgment.
- **Cons**:
  - Durability: weakest.
  - Message may be lost if the broker fails before it is written to disk.
  - No offset is not returned, so the producer cannot track the message.
- Use case: when you don't care about message durability and just want to send messages as fast as possible. Some metrics from IoT devices, Non-essential logs, etc.
- Records will still be replicated asynchronously.

```mermaid
sequenceDiagram
    participant P as 📤 Producer
    participant L as 🏢 Leader Broker
    participant F1 as 🔄 Follower 1
    participant F2 as 🔄 Follower 2
    
    Note over P,F2: ACKS = 0 (Fire and Forget)
    
    P->>L: 1. Send Message
    Note right of P: No waiting for ACK
    P->>P: 2. Continue immediately ⚡
    Note right of P: Fastest throughput
    
    L-->>F1: 3. Async replication
    L-->>F2: 4. Async replication
    
    Note over L,F2: Replication happens in background<br/>Risk: Message lost if leader fails
```

### ACKS = 1

- The producer waits for the leader to acknowledge the message.
- **Pros**:
  - Faster than ACKS = all.
  - Offset is returned, so the producer can track the message.
- **Cons**:
  - Durability: moderate.
  - If the leader fails after acknowledging the message but before it is replicated to followers, the message may be lost.
- **Use cases**: User activity logs, non-critical data that can be lost without significant impact.

```mermaid
sequenceDiagram
    participant P as 📤 Producer
    participant L as 🏢 Leader Broker
    participant F1 as 🔄 Follower 1
    participant F2 as 🔄 Follower 2
    
    Note over P,F2: ACKS = 1 (Leader ACK)
    
    P->>L: 1. Send Message
    L->>L: 2. Write to log
    L->>P: 3. ACK (offset=123) ✅
    Note right of P: Gets offset for tracking
    
    L-->>F1: 4. Async replication
    L-->>F2: 5. Async replication
    
    Note over L,F2: Risk: If leader fails before replication<br/>message could be lost
```


### ACKS = all (or -1)

- The producer waits for all in-sync replicas (ISRs) to acknowledge the message.
- Replication is synchronous.
- **Pros**:
  - Highest durability.
  - Guarantees that the message is written to all replicas before considering the request complete.
- **Cons**:
  - Slowest option, as it waits for all replicas to acknowledge.
  - Higher latency due to waiting for multiple acknowledgments.
- **Use cases**: Critical data that must be durable and available, such as financial transactions, order processing, etc.

```mermaid
sequenceDiagram
    participant P as 📤 Producer
    participant L as 🏢 Leader Broker
    participant F1 as 🔄 Follower 1
    participant F2 as 🔄 Follower 2
    
    Note over P,F2: ACKS = all (Maximum Durability)
    
    P->>L: 1. Send Message
    L->>L: 2. Write to log
    
    par Sync Replication
        L->>F1: 3. Replicate message
        F1->>L: 4. ACK from F1 ✅
    and
        L->>F2: 5. Replicate message
        F2->>L: 6. ACK from F2 ✅
    end
    
    L->>P: 7. Final ACK (offset=123) ✅
    Note right of P: All ISRs confirmed<br/>Highest durability
    
    Note over L,F2: All replicas have the message<br/>Safe even if leader fails
```

### Retries

- Types of errors that producer can encounter:
  - **Transient errors**: Temporary issues that can be retried, such as network failures, broker unavailability, etc.
    - LEADER_NOT_AVAILABLE
    - network errors
  - **Non-transient errors**: Permanent issues that cannot be retried, such as invalid messages, schema validation errors, etc.
    - INVALID_MESSAGE
    - IVALID_CONFIG
    - Too large message
- 2 types of retries:
  - **Automatic retries**: retry N times then give up.
  - **Manual retries**.
- Some metric can be used when retrying:
  - `retries`: Number of retries.
  - `retry.backoff.ms`: Time to wait before retrying.
  - `delivery.timeout.ms`: Maximum time to wait for a message to be delivered.
- **Issues**:
  - Duplicate messages: If a message is retried, it may be sent multiple times. This can lead to duplicate processing in the consumer.
  - Out of order messages: If a message is retried, it may be sent after a later message, leading to out-of-order processing in the consumer.

### Batching

> Kafka performance bottleneck is usually the network, not the disk. Batching is a way to improve throughput by sending multiple messages in a single request.

- Batching is the process of grouping multiple messages together and sending them in a single request. This method improves throughput by reducing the number of requests sent over the network.
- Batching is implemented in the producer.
- Control by 2 parameters:
  - `batch.size`: Maximum size of a batch in bytes. Default is 16KB.
  - `linger.ms`: Maximum time to wait before sending a batch. Default is 0 (send immediately).

#### `linger.ms`

- Producer will wait for up to the specified time before sending a batch.
- This add more latency, but allows more messages to be sent in a single request.
- The default value is 0, which means the producer will send messages immediately without waiting.
- If `linger.ms` = 0, batching still happens, but it is based on `batch.size` only. It may happen when the producer wait for the ACK from the broker.
- if `linger.ms` is set to a positive value, the producer will wait for that amount of time before sending a batch, even if the batch size is not reached.

> The more `linger.ms` is set, the more messages will be batched together, which leads to better throughput but higher latency.

#### `batch.size`

- Maximum size of a batch in bytes.
- Default is 16KB.
- If the batch size is reached, the producer will send the batch immediately, regardless of the `linger.ms` setting.
- A very large batch size can lead to increased memory usage and latency, as the producer will wait for more messages to fill the batch before sending.
- If a message is larger than the batch size, the the message won't be batched and will be sent immediately.

> The more `batch.size` is set, the more messages will be batched together, which leads to better throughput but higher latency.

#### Trade-offs

- Batching **improves throughput** by reducing the number of requests sent over the network but **adds latency.**
- May improve i/o performance by reducing the number of disk writes.
- Larger batch size can lead to increased memory usage and latency.

### Compression

- Compression is the process of reducing the size of messages before sending them over the network.
- Compression has a positive impact on performance:
  - Improve network throughput by reducing the amount of data sent over the network.
  - Save storage space on the broker.
  - Better latency due to reduced network traffic.
- Cons: CPU overhead for compressing and decompressing messages.
- Compression != Encryption. Compression is not a security feature, it is just a way to reduce the size of messages.
- **Types**: `none`, `gzip`, `snappy`, `lz4`, `zstd`.
- **Recommended**: `lz4` (normal case), `zstd` (high compression ratio).

| Compression Type | Compression Ratio | CPU Overhead | Compression Speed | Network bandwidth |
|------------------|-------------------|--------------|-------------------|-------------------|
| `gzip`             | Highest               | Highest          | Slowest            | Slowest            |
| `snappy`           | Medium                  | Medium              | Medium               | Medium               |
| `lz4`              | Slowest               | Slowest           | Fastest             | Highest             |
| `zstd`             | Medium                  | Medium              | Medium               | Medium               |

- **Compression is configured in the producer** and **decompressing is done in the consumer**.
- By default, brokers do not interfere with the batch when storing the batch.
- The compression format on the producer side must match the decompression format on the consumer side.

### Max In Flight Requests

- `max.in.flight.requests.per.connection`: Maximum number of unacknowledged requests the producer can send to a broker before blocking.
- Default is 5.
- Purpose: improve throughput by allowing multiple requests in producer (not improving throughput in the broker).

### Sticky Partitioner

> Awesome reference: [Kafka producer đã không còn Round Robin Partition với key null](https://thanhlv.com/blog/2024-08-07-Kafka-producer-da-khong-con-Round-Robin-Partition-voi-key-null.html)

- `Sticky Partitioner` is the logic of the producer that determines which partition a message should be sent to.
- Older versions of Kafka used round-robin partitioning, which mean producers don't immediately send messages but placing the partition-specific batches to sent later.
- In case the parameter `linger.ms` is set to a positive value, the producer will wait for that amount of time before sending a batch. In the enviroment with low throughput, this can lead to a situation where the producer is waiting for messages to fill the batches of each partition.

```mermaid
sequenceDiagram
    participant P as 📤 Producer
    participant P0 as 📂 Partition 0
    participant P1 as 📂 Partition 1
    participant P2 as 📂 Partition 2
    
    Note over P,P2: Round-Robin (Old Behavior)
    
    P->>P0: Msg1 → Batch 0 (1/16KB)
    P->>P1: Msg2 → Batch 1 (1/16KB)
    P->>P2: Msg3 → Batch 2 (1/16KB)
    P->>P0: Msg4 → Batch 0 (2/16KB)
    
    Note over P,P2: With linger.ms=10ms, producer waits<br/>for each partition batch to fill up
    
    P-->>P0: ⏰ Wait 10ms → Send partial batch
    P-->>P1: ⏰ Wait 10ms → Send partial batch  
    P-->>P2: ⏰ Wait 10ms → Send partial batch
    
    Note over P,P2: Result: 3 network requests<br/>with small, inefficient batches
```

- `Sticky Partitioner` improves the performance of the producer especially with high throughput. The producer sticky partitioner will:
  - Send messages to the same partition until the batch is full or the linger time is reached.
  - After that, it will switch to another partition.
  - This approach reduces the number of requests and improves throughput.
  - Sticky partitioner increases the rate of "filling" the batch, which leads to increase hit `batch.size` instead of `linger.ms`.
- We can say this is "Per Batch" Round-Robin, where the producer sends messages to the same partition until the batch is full, then switches to the next partition.

```mermaid
sequenceDiagram
    participant P as 📤 Producer
    participant P0 as 📂 Partition 0
    participant P1 as 📂 Partition 1
    participant P2 as 📂 Partition 2
    
    Note over P,P2: Sticky Partitioner (New Behavior)
    
    P->>P0: Msg1 → Batch 0 (1/16KB)
    P->>P0: Msg2 → Batch 0 (2/16KB)
    P->>P0: Msg3 → Batch 0 (3/16KB)
    P->>P0: Msg4 → Batch 0 (16KB - FULL!)
    
    P-->>P0: ✅ Send full batch immediately
    
    Note over P,P2: Switch to next partition
    
    P->>P1: Msg5 → Batch 1 (1/16KB)
    P->>P1: Msg6 → Batch 1 (2/16KB)
    
    Note over P,P2: Result: Fewer, larger batches<br/>Better network utilization
```

### Idempotence

> What if the producer sends the same message multiple times?

- `Idempotence`: the operation can be applied multiple times without changing the result beyond the initial application.
- How it works:
  - A Unique producer ID (PID) is assigned to each producer.
  - Each message sent by the producer includes a sequence number.
  - The broker checks the sequence number and PID to determine if the message is sent before.
- Configure:
  - `enable.idempotence=true`.
  - `acks=all`.
  - `retries=INT_MAX` (or a large number).
  - `max.in.flight.requests=5`
- Note: Idempotent producers only resolve ordering issues on producer side.

```mermaid
sequenceDiagram
    participant P as 📤 Producer (PID: 12345)
    participant B as 🏢 Broker
    participant L as 📝 Message Log
    
    Note over P,L: Idempotent Producer Flow
    
    %% First message attempt
    P->>B: 1. Send Message A (PID: 12345, Seq: 0)
    B->>B: 2. Check: PID 12345, Seq 0 - NEW
    B->>L: 3. Store Message A (PID: 12345, Seq: 0)
    B->>P: 4. ACK (offset: 100) ✅
    
    %% Second message
    P->>B: 5. Send Message B (PID: 12345, Seq: 1)
    B->>B: 6. Check: PID 12345, Seq 1 - NEW
    B->>L: 7. Store Message B (PID: 12345, Seq: 1)
    B->>P: 8. ACK (offset: 101) ✅
    
    %% Network failure scenario - retry
    Note over P,B: Network timeout, producer retries Message B
    
    P->>B: 9. Retry Message B (PID: 12345, Seq: 1)
    B->>B: 10. Check: PID 12345, Seq 1 - DUPLICATE!
    Note right of B: Already processed this sequence
    B->>P: 11. ACK (offset: 101) ✅
    Note right of B: Returns same offset,<br/>no duplicate in log
    
    %% Next message continues normally
    P->>B: 12. Send Message C (PID: 12345, Seq: 2)
    B->>B: 13. Check: PID 12345, Seq 2 - NEW
    B->>L: 14. Store Message C (PID: 12345, Seq: 2)
    B->>P: 15. ACK (offset: 102) ✅
```

### Serialization

> Serialization is the process of converting objects or data structures into bytes for transmission or storage.

#### Formats

- Text-based formats:
  - `JSON`.
  - `XML`.
  - `String`.
- Binary-based formats:
  - `Avro`.
  - `Protobuf`.
  - `Thrift`.
  
#### Format Selection

- **Complexity**: The difficulty of parsing and generating the format.
  - Text-based formats are easier to read and debug.
  - Binary formats are more efficient in terms of size and speed.
  - Binary formats are more suitable for high-performance applications.
- **Compatibility**: Is the format compatible with other components in the system?
  - Text-based formats are more flexible and easier to evolve.
  - Binary formats require careful schema management.
- **Size**: The size of the serialized data.
  - Binary formats are usually smaller than text-based formats.
  - Smaller size leads to better network throughput and storage efficiency.

## Consumer

- A consumer is an application that subscribes to Kafka topics and processes messages.
- Consumer can read from:
  - The beginning of the topic.
  - The specified offset or timestamp using `seek()`.
  - The current position: a consumer goes down then restart or a new consumer joins the group and need to know the latest position.

### Offset commit

- Consumers need to keep track of the offsets of the messages they have processed.
- Offset commit is the process of saving the current position of the consumer in a partition of the topic.
- Where to store commits?
  - Offsets are stored in a special topic called `__consumer_offsets`. This is the default and recommended way.
  - Current offset = committed offset + 1.

### Delivery Guarantees

- Dekivery Guarantee: The assurance provided by the system about message delivery reliability and consistency.
- 3 types of delivery guarantees:
  - **At-Most-Once**: A message us processed at most once but it may be error.
    - If a message is lost, it will not be retried.
    - ***Usecase***: Logging, metrics, event IoT, etc.
  - **At-Least-Once**: A message is processed at least once, but it may be duplicated.
    - If a message is lost, it will be retried.
    - ***Usecase***: Email, SMS, Data replication, etc.
  - **Exactly-Once**: A message is processed exactly once, without duplicates or loss.
    - This is the most complex and expensive guarantee to achieve.
    - ***Usecase***: Financial transactions, where duplicates or loss are not acceptable.
    - **Note**: A messaging platform cannot offer exactly-once delivery guarantee by itself. It requires the application to be designed to handle idempotent operations and deduplication logic.

> Awesome reference: [You Cannot Have Exactly-Once Delivery](https://bravenewgeek.com/you-cannot-have-exactly-once-delivery/)
> In the letter I mail you, I ask you to call me once you receive it. You never do. Either you really didn’t care for my letter or it got lost in the mail. That’s the cost of doing business. I can send the one letter and hope you get it, or I can send 10 letters and assume you’ll get at least one of them. The trade-off here is quite clear (postage is expensive!), but sending 10 letters doesn’t really provide any additional guarantees. In a distributed system, we try to guarantee the delivery of a message by waiting for an acknowledgement that it was received, but all sorts of things can go wrong. Did the message get dropped? Did the ack get dropped? Did the receiver crash? Are they just slow? Is the network slow? Am I slow? FLP and the Two Generals Problem are not design complexities, they are impossibility results.
> To reiterate, there is no such thing as exactly-once delivery. We must choose between the lesser of two evils, which is at-least-once delivery in most cases. This can be used to simulate exactly-once semantics by ensuring idempotency or otherwise eliminating side effects from operations. Once again, it’s important to understand the trade-offs involved when designing distributed systems. There is asynchrony abound, which means you cannot expect synchronous, guaranteed behavior. Design for failure and resiliency against this asynchronous nature.
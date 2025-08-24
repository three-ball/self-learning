# HASH

![Hash](../images/hash.png)

## Definition

- **A Hash Table** (also called Hash Map) is a data structure that implements an associative array abstract data type, a structure that can map keys to values. It uses a hash function to compute an index into an array of buckets or slots, from which the desired value can be found. The main goal is to achieve O(1) average time complexity for search, insertion, and deletion operations.

## Characteristics

- **Key-Value Pair**: Hash tables store data in key-value pairs, allowing for efficient retrieval based on the key.
- **Hash Function**: A hash function is used to compute an index from the key, which determines where the value is stored in the array.
- **Collision Handling**: When two keys hash to the same index, a collision occurs. Common strategies for handling collisions include chaining (using linked lists) and open addressing (finding another open slot).
- **Dynamic Resizing**: Many hash table implementations resize the underlying array when the load factor (number of elements / array size) exceeds a certain threshold, maintaining efficient performance.

## Core concepts

```mermaid
graph TD
    A[Start Operation] --> B{Operation Type?}
    
    B -->|PUT| C[Calculate Hash of Key]
    C --> D[Find Bucket Index]
    D --> E{Bucket Empty?}
    E -->|Yes| F[Create New Node]
    F --> G[Increment Count]
    G --> H{Load Factor > 0.75?}
    H -->|Yes| I[Resize & Rehash]
    H -->|No| Z[End]
    I --> Z
    E -->|No| J[Traverse Chain]
    J --> K{Key Exists?}
    K -->|Yes| L[Update Value]
    L --> Z
    K -->|No| M[Add to End of Chain]
    M --> G
    
    B -->|GET| N[Calculate Hash of Key]
    N --> O[Find Bucket Index]
    O --> P[Traverse Chain]
    P --> Q{Key Found?}
    Q -->|Yes| R[Return Value]
    Q -->|No| S[Return Not Found]
    R --> Z
    S --> Z
    
    B -->|DELETE| T[Calculate Hash of Key]
    T --> U[Find Bucket Index]
    U --> V[Search & Remove Node]
    V --> W[Decrement Count]
    W --> Z
    
    style F fill:#e8f5e8
    style L fill:#fff2cc
    style R fill:#e1f5fe
    style I fill:#fce4ec
```

- **Hash function**: a function that takes a key (input) and returns an index in the hash table array.

![Hash Table](../images/hash_tb.png)

- **Collision**: When two different keys produce the same hash value.
- **Chaining**: A collision resolution technique where each bucket in the hash table contains a linked list of all entries that hash to the same index.

![Chaining](../images/chaining.png)

- **Load factor**: The load factor is basically how much data we have in our table. This load factor is called `lambda`. This factor is going to vary depending on how our data structure resolves collisions. We need to resize table when `lambda` exceeds a certain threshold to maintain efficient performance.

```sh
lambda = number of entries / size of the hash table
```
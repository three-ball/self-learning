# MIT 6.824 Distributed System

> This is a note for the MIT 6.824 Distributed System course in 2020 - Spring, a free course on YouTube.

- [MIT 6.824 Distributed System](#mit-6824-distributed-system)
  - [Lecture 1: Introduction](#lecture-1-introduction)
    - [What is a Distributed System?](#what-is-a-distributed-system)
    - [Case study: MapReduce (MR)](#case-study-mapreduce-mr)

## Lecture 1: Introduction

- Link: [Lecture 1: Introduction](https://youtu.be/cQP8WApzIQQ?si=p-Zn1kFrTYjteVz7)

### What is a Distributed System?
- **Distributed System:** A group of computers cooperating to provide a service.
- It's **not easy to build distributed systems**:
  - Concurrency
  - Complex Infras
  - How to get high Perf?
  - Parital failure.
- Main topics in this course:
  - `Storage`
  - `Communication`
  - `Computation`
- Distributed features:
  - `Fault Tolerance`: What if something broken?
    - *Availability*
    - *Replication*
    - *Recoverability*
  - `Consistency`: Genral purpose infrastructure needs to be well-defined behavior, something like: read(x) value from the most recent write(x) on replica...
  - `Performace`: 2x Resource will bring to 2x Perf or 2x throughput?
  - `Trade-off`: Fault tolerance, consistency and performance ara enemies.
    - Fault tolerance and consistency require communication but communication is often slow!
  - `Implementation`: RPC, Thread, Concurency...

### Case study: MapReduce (MR)

> **MapReduce** is a programming model and an associated implementation for processing and generating large data sets. Users specify a map function that processes a key/value pair to generate a set of intermediate key/value pairs, and a reduce function that merges all intermediate values associated with the same intermediate key.

![alt text](images/lect_01/mr_overview.png)

- **Useful link**
  - [MapReduce: Simplified Data Processing on Large Clusters](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)
  - [Note](https://pdos.csail.mit.edu/6.824/notes/l01.txt)

![alt text](images/lect_01/mr_overviews.png)

- MR scales well because:
  - `Nx <worker>` computer might get `Nx <throughput>`. Map() can run in parallel, same for Reduce().
  - More workers, more throughput.
- MR rules:
  - **Map()** and **Reduce()** functions are pure functions.
  - **Map()** and **Reduce()** functions are stateless and not interact with other workers.
- In 2004, Google published the paper about MR, and it's the first paper that introduces the concept of distributed systems. The bottleneck of MR is network speed.
- How MR minimize network traffic:
  - Coordinater tries to run each Map task on GFS (Google File System) node that contains the input data. This makes the data local to the computation.
  - Map Workers write intermediate data to local disk.
  - Reduce Workers read intermediate data from map workers over the network.
- Suppose MR runs a Map task twice, one Reduce sees first run's output and the other sees the second run's output. This is a consistency problem. The two Map executions had better produce identical intermediate output! **They are only allowed to look at their arguments/input, no state, no file I/O, no interaction, no external communication, no random numbers.**
- Current status of MR: `Spark`, `Flink`, `Hadoop`, etc. No longer in use at Google. Google uses a system called `FlumeJava` and GFS is replaced by `Colossus`.
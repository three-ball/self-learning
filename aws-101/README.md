# AWS 101

> This repository contains notes about AWS - from zero.

## "Compute in the Cloud"

### EC2 - Elastic Compute Cloud

- **Amazon EC2**: A service providing virtual servers to power your business applications.
- EC2 Instance Types:
    - `General purpose instances`: Balance of compute, memory, and networking resources.
        - Usecases: Web services, code repositories, development environments.
    - `Compute optimized instances`: Ideal for compute-bound applications that benefit from high-performance processors.
        - Usecases: High-performance web servers, scientific modeling, dedicated gaming servers.
    - `Memory optimized instances`: Designed to deliver fast performance for workloads that process large data sets in memory.
        - Usecases: Databases, real-time big data analytics, in-memory caching.
    - `Accelerated computing instances`: Use hardware accelerators, or co-processors, to perform functions such as floating-point number calculations, graphics processing, or data pattern matching more efficiently than software running on a general-purpose CPU.
        - Usecases: Graphics processing, floating-point calculations, data pattern matching.
    - `Storage optimized instances`: Designed for workloads that require high, sequential read and write access to very large data sets on local storage. They are optimized to deliver tens of thousands of low-latency, random I/O operations per second (IOPS) to applications.
        - Usecases: Data warehousing applications, distributed file systems, network file systems, log or data processing applications.

- [EC2 Naming Conventions](https://docs.aws.amazon.com/ec2/latest/instancetypes/instance-type-names.html): Instance types are named based on their instance family and instance size. Example: `t2.micro`.
    - `t2.micro`: `t` stands for `instance type`, `2` stands for `generation`, `micro` stands for `size`.

![alt text](images/naming_convention.png)

- EC2 Pricing:
    - **On-Demand Instances**: Pay per hour or per second depending on which instances you run and OS. No long-term commitments or upfront payments.
        - Ideal for short-term, irregular workloads that cannot be interrupted. No upfront costs or minimum contracts apply.
    - **Savings Plans**: Commit to a consistent usage (USD/hour) for 1 or 3 years. Save up to 72% on EC2, Fargate, and Lambda.
        - Ideal for workloads with *steady state usage*, predictable usage, long-term savings.
    - **Reserved Instances**: Suited for steady-state or predictable workloads. Up to 75% discount compared to On-Demand pricing. Commit to a 1 or 3-year term. Can be: All Upfront, Partial Upfront, No Upfront.
        - Ideal for *steady workloads* with upfront or partial payment options.
    - **Spot Instances**: Request unused EC2 instances at steep discounts. Can be terminated by AWS with 2 minutes warning. Up to 90% discount compared to On-Demand pricing.
        - Ideal for workloads that can tolerate interruptions (e.g., batch processing). Workloads should be fault-tolerant and flexible. Should has some gracefull shutdown mechanism.
    - **Dedicated Hosts**: Most expensive! Physical EC2 server dedicated for your use. Meet compliance requirements by ensuring no shared tenancy.
        - Ideal for Compliance-driven, fully isolated resources.

### EC2 Auto Scaling

![alt text](images/scaling_demand.png)

- **Amazon EC2 Auto Scaling** enables you to automatically add or remove Amazon EC2 instances in response to changing application demand. There are two approaches to scaling:
    - **Manual Scaling**: responds to changing demand. 
    - **Predictive  Scaling**: automatically schedules the right number of Amazon EC2 instances based on predicted demand.
- There 03 configured parameters for Auto Scaling:
    - **Minimum Capacity**: The minimum number of instances that Auto Scaling maintains.
    - **Desired Capacity**: The number of instances that Auto Scaling tries to maintain.
    - **Maximum Capacity**: The maximum number of instances that Auto Scaling maintains.

![alt text](images/auto_scaling_group.png)

- Like the image above:
    - Minimum Capacity: Ensures at least 1 instance is always running, providing continuous availability.
    - Desired Capacity: Allows the system to scale up to 2 instances during normal operations for optimal performance.
    - Maximum Capacity: Ensures that the system can scale up to 4 instances during peak demand to maintain

### Elastic Load Balancing (ELB)

- **Elastic Load Balancing** is the AWS service that automatically distributes incoming application traffic across multiple resources, such as Amazon EC2 instances.
    - AWS101 Question: What if the ELB is down? How to handle this situation? Answer: [StackOverflow](https://stackoverflow.com/questions/46698011/are-amazon-elastic-load-balancer-elb-failure-proof)
- How ELB works (with EC2 Auto Scaling):
    - Auto-scaling services notify ELB when new instances are ready.
    - ELB directs traffic to new instances as they come online.
    - During scale-down, ELB waits for requests to complete before terminating instances.

### Amazon Simple Notification Service (SNS) & Amazon Simple Queue Service (SQS)

- **Amazon Simple Notification Service (Amazon SNS)** is a publish/subscribe service. Using Amazon SNS topics, a publisher publishes messages to subscribers. 
- **Amazon Simple Queue Service (Amazon SQS)** is a message queuing service. It allows you to decouple and scale microservices, distributed systems, and serverless applications.

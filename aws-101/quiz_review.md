# AWS Cloud Practitioner Essentials Quiz Review

## [AWS Skills Centers: Becoming a Cloud Practitioner - Mini-Quiz Review](https://explore.skillbuilder.aws/learn/courses/14704/aws-skills-centers-becoming-a-cloud-practitioner-mini-quiz-review)

1. Which of the following are `regional AWS services`? (Choose two.)
    - [ ] Amazon Shield -> available globally
    - [x] Amazon Transit Gateway -> regional resource
    - [x] Amazon S3 -> regional resource
    - [ ] Amazon Route 53 -> available globally
    - [ ] Amazon Organizations -> available globally

- [List of AWS Global Services](https://stackoverflow.com/questions/68811957/aws-global-services):
    - IAM
    - CloudFront
    - Route 53
    - WAF
- Most AWS Services are regional services, For example:
    - EC2
    - Beanstalk
    - Lambda
    - Rekognition
- AWS Organizations is global service [AWS Organizations endpoints and quotas](https://docs.aws.amazon.com/general/latest/gr/ao.html#:~:text=Because%20AWS%20Organizations%20is%20a,AWS%20Regions%20in%20each%20partition.)
- AWS Shield is global service [AWS Shield FAQs](https://aws.amazon.com/shield/faqs/#:~:text=AWS%20Shield%20Advanced%20is%20available,in%20front%20of%20your%20application.)

2. You are learning more about securing AWS resources. Which of the following allows to **group user account together and assign permissions** to those groups?
    - [ ] Resource groups
    - [ ] AWS Organizations
    - [x] AWS IAM
    - [ ] Tagging

3. The DevOps team at an IT company is moving 500 GB of data from EC2 to S3 bucket in the same region. Which of the following scenarios captures the correct charges for this data transfer?
    - [ ] The company would only be charged for the inbound data transfer into the S3 bucket
    - [ ] The company would only be charged for the outbound data transfer from the EC2 instance
    - [ ] The company would be charged for both the inbound and outbound data transfer
    - [x] The company would not be charged for the data transfer

The company would not be charged for this data transfer. There are three fundamental drivers of cost with AWS: compute, storage, and outbound data transfer. In most cases, there is no charge for inbound data transfer or data transfer between other AWS services within the same region.

4. Which of the following is an AWS database service?
    - [x] Amazon Redshift
    - [ ] Amazon CloudFront
    - [ ] Amazon S3
    - [ ] Amazon Route 53

Amazon database services:

| Service Name       | Type                     |
|--------------------|--------------------------|
| Amazon RDS         | Relational Database      |
| Amazon Redshift    | Relational Database      |
| Amazon Aurora      | Relational Database      |
| Amazon DynamoDB    | NoSQL Database           |
| Amazon DocumentDB  | Document Database        |
| Amazon Neptune     | Graph Database           |
| Amazon Timestream  | Time Series Database     |
| Amazon Keyspaces   | Managed Cassandra Service|

5. Which AWS services can be used to facilitate organizational change management, part of the Reliability pillar of AWS Well-Architected Framework? 
    [ ] AWS Trust Advisor
    [x] AWS Config
    [x] AWS CloudTrail
    [x] AWS CloudWatch

- There are three best practices for AWS Reliability pillar:
    - Change Management
    - Failure Management
    - Demand Management

6. Which of the following are correct statements regarding the AWS Shared Responsibility Model? (Select two)
    - [x] For abstracted services like Amazon S3, AWS operates the infrastructure layer, the operating system, and platforms
    - [x] AWS is responsible for Security 'of' the Cloud
    - [ ] For a service like Amazon EC2, that falls under Infrastructure as a Service (IaaS), AWS is responsible for maintaining guest operating system
    - [ ] AWS is responsible for training AWS and customer employees on AWS products and services
    - [ ] Configuration Management is the responsibility of the customer

**For abstracted services like Amazon S3, AWS operates the infrastructure layer, the operating system, and platforms:**
For abstracted services, such as Amazon S3 and Amazon DynamoDB, AWS operates the infrastructure layer, the operating system, and platforms, and customers access the endpoints to store and retrieve data.
**AWS is responsible for Security 'of' the Cloud**
AWS is responsible for protecting the infrastructure that runs all of the services offered in the AWS Cloud. This infrastructure is composed of the hardware, software, networking, and facilities that run AWS Cloud services.
**Configuration Management is the responsibility of the customer**
Configuration management is a shared responsibility. AWS maintains the configuration of its infrastructure devices, but a customer is responsible for configuring their own guest operating systems, databases, and applications.

7. Which benefit of Cloud Computing allows AWS to offer lower pay-as-you-go prices as usage from hundreds of thousands of customers is aggregated in the cloud?
    - [ ] Trade capital expense for variable expense
    - [ ] Increased speed and agility
    - [ ] Go global in minutes
    - [x] Massive economies of scale

**Massive economies of scale**
Because usage from hundreds of thousands of customers is aggregated in the cloud, providers such as AWS can **achieve higher economies of scale, which translates into lower pay-as-you-go** prices.

8. Which of the following AWS services can be used to forecast your AWS account usage and costs?
    - [ ] AWS Budgets
    - [ ] AWS Config
    - [ ] AWS Cost & Usage Report (AWS CUR)
    - [x] AWS Cost Explorer

**AWS Cost Explorer**
AWS Cost Explorer is a free tool that allows you to view your costs and usage over time and forecast how much you are likely to spend in the future. 

9. Due to regulatory and compliance reasons, an organization is supposed to use a hardware device for any data encryption operations in the cloud. Which AWS service can be used to meet this compliance requirement?
    - [ ] AWS Key Management Service (KMS)
    - [ ] AWS Certificate Manager
    - [x] AWS CloudHSM
    - [ ] AWS Secrets Manager

**AWS CloudHSM**
AWS CloudHSM is a cloud-based Hardware Security Module (HSM) that enables you to easily generate and use your encryption keys on the AWS Cloud.

10. Which AWS service can help you analyze your infrastructure to identify unattached or underutilized Amazon EBS Elastic Volumes?
    - [ ] AWS Config
    - [x] AWS Trusted Advisor
    - [ ] AWS CloudTrail
    - [ ] AWS Cost Explorer

**AWS Trusted Advisor**
AWS Trusted Advisor is an online tool that provides you real-time guidance to help you provision your resources following AWS best practices. Trusted Advisor can help you analyze your infrastructure to identify unattached or underutilized Amazon EBS Elastic Volumes.

11. A developer would like to automate operations on his on-premises environment using Chef and Puppet. Which AWS service can help with this task?
    - [ ] AWS CloudFormation
    - [x] AWS OpsWorks
    - [ ] AWS Systems Manager
    - [ ] AWS CodeDeploy

**AWS OpsWorks**
AWS OpsWorks is a configuration management service that helps you configure and operate applications of all shapes and sizes using Chef and Puppet.

12. An e-commerce company wants to store data from a recommendation engine in a database. As a Cloud Practioner, which AWS service would you recommend to provide this functionality with the LEAST operational overhead for any scale?
    - [ ] Amazon RDS
    - [x] Amazon DynamoDB
    - [ ] Amazon Redshift
    - [ ] Amazon Neptune

**Amazon DynamoDB**
Amazon DynamoDB is a key-value and document database that delivers sub-millisecond performance at any scale. You can use Amazon DynamoDB to store recommendation results with the LEAST operational overhead for any scale.
**Amazon Neptune**
Amazon Neptune is a fully managed database service built for the cloud that makes it easier to build and run graph applications. It's not the right fit to store recommendation results with the LEAST operational overhead for any scale.

13. An IT company wants to run a log backup process every Monday at 2 AM. The usual runtime of the process is 5 minutes. As a Cloud Practitioner, which AWS services would you recommend to build a serverless solution for this use-case? (Select two)
    - [x] AWS Lambda
    - [x] Amazon Eventbridge
    - [ ] Amazon Step Function
    - [ ] Amazon S3

**AWS Lambda**
AWS Lambda lets you run code without provisioning or managing servers. You pay only for the compute time you consume. 
**Amazon Eventbridge**
Amazon EventBridge is a service that provides real-time access to changes in data in AWS services, your own applications, and software as a service (SaaS) applications without writing code. Amazon EventBridge Scheduler is a serverless task scheduler that simplifies creating, executing, and managing millions of schedules across AWS services without provisioning or managing underlying infrastructure.

14. A company's flagship application runs on a fleet of Amazon Elastic Compute Cloud (Amazon EC2) instances. As per the new policies, the system administrators are looking for the best way to provide secure shell access to Amazon Elastic Compute Cloud (Amazon EC2) instances without opening new ports or using public IP addresses.

Which tool/service will help you achieve this requirement?

    - [ ] AWS Systems Manager
    - [ ] AWS Config
    - [ ] AWS CloudTrail
    - [x] AWS Systems Manager Session Manager

**AWS Systems Manager Session Manager**
AWS Systems Manager Session Manager is a fully-managed service that provides you with an interactive browser-based shell and CLI experience. It helps provide secure and auditable instance management without the need to open inbound ports, maintain bastion hosts, and manage SSH keys. 

15. Which service gives a personalized view of the status of the AWS services that are part of your Cloud architecture so that you can quickly assess the impact on your business when AWS service(s) are experiencing issues?

    - [ ] AWS Config
    - [ ] AWS CloudTrail
    - [ ] AWS Trusted Advisor
    - [x] AWS Personal Health Dashboard

**AWS Personal Health Dashboard**
AWS Health - Your Account Health Dashboard provides alerts and remediation guidance when AWS is experiencing events that may impact you.

16. AWS Organizations provides which of the following benefits? (Select two)

    - [x] Share the reserved Amazon EC2 instances amongst the member AWS accounts
    - [x] Volume discounts for Amazon EC2 and Amazon S3 aggregated across the member AWS accounts
    - [ ] Centralized logging
    - [ ] Centralized monitoring
    - [ ] Centralized security

**Share the reserved Amazon EC2 instances amongst the member AWS accounts**
AWS Organizations allows you to share reserved Amazon EC2 instances amongst the member AWS accounts.
**Volume discounts for Amazon EC2 and Amazon S3 aggregated across the member AWS accounts**
AWS Organizations allows you to aggregate volume discounts for Amazon EC2 and Amazon S3 across the member AWS accounts.

17. Which of the following AWS services have data encryption automatically enabled? (Select two)?

    - [x] Amazon S3
    - [x] AWS Storage Gateway
    - [ ] Amazon Elastic Block Store (Amazon EBS)
    - [ ] Amazon Elastic File System (Amazon EFS)

**Amazon S3**
All Amazon S3 buckets have encryption configured by default, and objects are automatically encrypted by using server-side encryption with Amazon S3 managed keys (SSE-S3).
**AWS Storage Gateway**
AWS Storage Gateway is a hybrid cloud storage service that gives you on-premises access to virtually unlimited cloud storage. All data transferred between the gateway and AWS storage is encrypted using SSL.

18. According to the AWS Cloud Adoption Framework (AWS CAF), what are two tasks that a company should perform when planning to migrate to the AWS Cloud and aiming to become more responsive to customer inquiries and feedback as part of their organizational transformation?

    - [x] Organize your teams around products and value streams
    - [x] Leverage agile methods to rapidly iterate and evolve

19. Which of the following AWS services can be used to connect a company's on-premises environment to a VPC without using the public internet?

    - [x] AWS Direct Connect
    - [ ] Amazon VPC
    - [ ] Amazon Route 53
    - [ ] AWS Site-to-Site VPN

**AWS Direct Connect**
You can use AWS Direct Connect to establish a private virtual interface from your on-premise network directly to your Amazon VPC, providing you with a private, high bandwidth network connection between your network and your VPC. **This connection is private and does not go over the public internet**. It takes at least a month to establish this physical connection.

**AWS Site-to-Site VPN** - AWS Site-to-Site VPN creates a secure connection between your data center or branch office and your AWS cloud resources. This connection goes over the public internet.

20. Which of the following are correct statements regarding the AWS Global Infrastructure? (Select two)

    - [ ] Each Availability Zone (AZ) consists of two or more discrete data centers
    - [ ] Each AWS Region consists of a minimum of two Availability Zones (AZ)
    - [x] Each Availability Zone (AZ) consists of one or more discrete data centers
    - [x] Each AWS Region consists of a minimum of three Availability Zones (AZ)
    - [ ] Each AWS Region consists of two or more Edge Locations

**Each Availability Zone (AZ) consists of one or more discrete data centers**
Each Availability Zone (AZ) consists of one or more discrete data centers, each with redundant power, networking, and connectivity, housed in separate facilities.
**Each AWS Region consists of a minimum of three Availability Zones (AZ)**
Each AWS Region consists of a minimum of three Availability Zones (AZ) that are geographically separated from each other within the same geographic region.

21. Which of the following AWS services support reservations to optimize costs? (Select three)

    - [x] Amazon EC2
    - [x] Amazon RDS
    - [x] Amazon Redshift
    - [ ] Amazon S3
    - [ ] Amazon Route 53

**AWS services support reservations**:
- Amazon EC2 Reserved Instances (RI)
- Amazon RDS RIs
- Amazon DynamoDB Reserved Capacity
- Amazon ElastiCache Reserved Nodes
- Amazon Redshift Reserved Nodes





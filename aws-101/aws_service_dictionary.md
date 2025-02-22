# AWS SERVICE DICTIONARY

## Table of Contents

- [Services/Keywords](#serviceskeywords)
  - [1. AWS CloudHSM](#1-aws-cloudhsm)
  - [2. AWS Key Management Service (KMS)](#2-aws-key-management-service-kms)
  - [3. AWS Secrets Manager](#3-aws-secrets-manager)
  - [4. AWS Cost Explorer](#4-aws-cost-explorer)
  - [5. AWS Pricing Calculator](#5-aws-pricing-calculator)
  - [6. AWS Budgets](#6-aws-budgets)
  - [7. AWS Trusted Advisor](#7-aws-trusted-advisor)
  - [8. AWS Inspector](#8-aws-inspector)
  - [9. AWS Config](#9-aws-config)
  - [10. AWS OpsWorks](#10-aws-opsworks)
  - [11. AWS CloudFormation](#11-aws-cloudformation)
  - [12. AWS Lambda](#12-aws-lambda)
  - [13. AWS EventBridge](#13-aws-eventbridge)
  - [14. AWS Step Functions](#14-aws-step-functions)
  - [15. AWS IAM Identity Center](#15-aws-iam-identity-center)
  - [16. AWS IAM](#16-aws-iam)
  - [17. AWS Cognito](#17-aws-cognito)
  - [18. AWS Systems Manager Session Manager](#18-aws-systems-manager-session-manager)
  - [19. EC2 Instance Connect](#19-ec2-instance-connect)
  - [20. AWS Health - Your Account Health Dashboard](#20-aws-health---your-account-health-dashboard)
  - [21. AWS Health - Service Health Dashboard](#21-aws-health---service-health-dashboard)
  - [22. AWS Systems Manager](#22-aws-systems-manager)
  - [23. Amazon DynamoDB with global tables](#23-amazon-dynamodb-with-global-tables)
  - [24. AWS Local Zones](#24-aws-local-zones)
  - [25. AWS Edge Locations](#25-aws-edge-locations)

## Services/Keywords

### 1. AWS CloudHSM
- **Keyword**: `Security`, `Encryption`, `Compliance`, `Hardware`.
- **Description**: AWS CloudHSM is a cloud-based hardware security module (HSM) that enables you to easily generate and use your own encryption keys on the AWS Cloud.

### 2. AWS Key Management Service (KMS)
- **Keyword**: `Security`, `Encryption`, `Compliance`.
- **Description**: AWS Key Management Service (AWS KMS) makes it easy for you to create and manage cryptographic keys and control their use across a wide range of AWS services and in your applications.

### 3. AWS Secrets Manager
- **Keyword**: `Security`, `Credential`.
- **Description**: AWS Secrets Manager helps you protect secrets needed to access your applications, services, and IT resources.

### 4. AWS Cost Explorer
- **Keyword**: `Cost Management`, `Forecasting`.
- **Description**: AWS Cost Explorer is a free tool that allows you to view your costs and usage over time and forecast how much you are likely to spend in the future.

### 5. AWS Pricing Calculator
- **Keyword**: `Cost Management`.
- **Description**: The AWS Pricing Calculator is a tool that allows you to estimate the cost of using AWS services based on your specific needs.

### 6. AWS Budgets
- **Keyword**: `Cost Management`, `Arlert`.
- **Description**:  AWS Budgets gives the ability to set custom budgets that alert you when your costs or usage exceed (or are forecasted to exceed) your budgeted amount.  *AWS Budgets cannot forecast your AWS account cost and usage*.

### 7. AWS Trusted Advisor
- **Keyword**: `Cost Management`, `Performance`, `Security`, `Fault Tolerance`, `Service Limits`.
- **Description**: AWS Trusted Advisor is an online tool that provides you real-time guidance to help you provision your resources following AWS best practices.

### 8. AWS Inspector
- **Keyword**: `Security`, `Vulnerability`.
- **Description**:  Amazon Inspector is an automated security assessment service that helps improve the security and compliance of applications deployed on your Amazon EC2 instances.

### 9. AWS Config
- **Keyword**: `Compliance`, `Configuration Management`, `Change Management`.
- **Description**: AWS Config is a service that enables you to assess, audit, and evaluate the configurations of your AWS resources.

### 10. AWS OpsWorks
- **Keyword**: `Automation`, `Chef`, `Puppet`.
- **Description**: AWS OpsWorks is a configuration management service that helps you configure and operate applications of all shapes and sizes using Chef and Puppet.

### 11. AWS CloudFormation
- **Keyword**: `Automation`.
- **Description**: AWS CloudFormation provides a common language for you to describe and provision all the infrastructure resources in your cloud environment.

### 12. AWS Lambda
- **Keyword**: `Serverless`, `Compute`.
- **Description**: AWS Lambda lets you run code without provisioning or managing servers. You pay only for the compute time you consume.

### 13. AWS EventBridge
- **Keyword**: `Serverless`, `Event-Driven`.
- **Description**: Amazon EventBridge is a serverless event bus service that makes it easy to connect applications together using data from your own applications, integrated Software as a Service (SaaS) applications, and AWS services.  Amazon EventBridge Scheduler is a serverless task scheduler that simplifies creating, executing, and managing millions of schedules across AWS services without provisioning or managing underlying infrastructure.

### 14. AWS Step Functions
- **Keyword**: `Serverless`, `Workflow`.
- **Description**: AWS Step Function lets you **coordinate multiple AWS services into serverless workflows**. You can design and run workflows that stitch together services such as AWS Lambda, AWS Glue and Amazon SageMaker. Step Function cannot be used to run a process on a schedule.

### 15. AWS IAM Identity Center
- **Keyword**: `Security`, `Identity`, `Access Management`, `SSO`.
- **Description**: AWS IAM Identity Center is the successor to AWS Single Sign-On (AWS SSO). It is built on top of AWS Identity and Access Management (IAM) to simplify access management to multiple AWS accounts, AWS applications, and other SAML-enabled cloud applications.

### 16. AWS IAM
- **Keyword**: `Security`, `Identity`, `Access Management`.
- **Description**:  AWS Identity and Access Management (AWS IAM) enables you to securely control access to AWS services and resources for your users. It is not used to log in but to manage users and roles.

### 17. AWS Cognito
- **Keyword**: `Security`, `Identity`, `Authentication`, `Authorization`.
- **Description**: Amazon Cognito lets you add user sign-up, sign-in, and access control to your web and mobile apps quickly and easily. Amazon Cognito scales to millions of users and supports sign-in with social identity providers, such as Apple, Facebook, Google, and Amazon, and enterprise identity providers via SAML 2.0.

### 18. AWS Systems Manager Session Manager
- **Keyword**: `Management`, `Remote Access`, `No SSH`.
- **Description**: AWS Systems Manager Session Manager is a fully managed AWS Systems Manager capability that lets you manage your Amazon EC2 instances through an interactive one-click browser-based shell or through the AWS CLI. Session Manager provides secure and auditable instance management without the need to open inbound ports, maintain bastion hosts, or manage SSH keys.

### 19. EC2 Instance Connect
- **Keyword**: `Management`, `Remote Access`, `SSH`.
- **Description**: EC2 Instance Connect provides a simple and secure way to connect to your instances using Secure Shell (SSH). With EC2 Instance Connect, you can control SSH access to your instances using AWS Identity and Access Management (IAM) policies as well as audit connection requests with AWS CloudTrail events. EC2 Instance Connect will need port 22 to be open for traffic.

### 20. AWS Health - Your Account Health Dashboard
- **Keyword**: `Management`, `Health`.
- **Description**: AWS Health - Your Account Health Dashboard provides alerts and remediation guidance when AWS is experiencing events that may impact you.

### 21. AWS Health - Service Health Dashboard
- **Keyword**: `Management`, `Health`.
- **Description**: The AWS Health - Service Health Dashboard is the single place to learn about the availability and operations of AWS services. You can view the overall status of AWS services, and you can sign in to view personalized communications about your particular AWS account or organization.

### 22. AWS Systems Manager
- **Keyword**: `Management`, `Automation`, `Patch Management`.
- **Description**: AWS Systems Manager allows you to centralize operational data from multiple AWS services and automate tasks across your AWS resources. You can create logical groups of resources such as applications, different layers of an application stack, or production versus development environments.

### 23. Amazon DynamoDB with global tables
- **Keyword**: `Database`, `NoSQL`, `Global`, `Active-Active`, `Cross-Region`.
- **Description**: Amazon DynamoDB global tables provide a fully managed solution for deploying a multi-region, multi-master database, without having to build and maintain your own replication solution.

### 24. AWS Local Zones
- **Keyword**: `Compute`, `Low Latency`.
- **Description**: AWS Local Zones allow you to use select AWS services, like compute and storage services, closer to more end-users, providing them very low latency access to the applications running locally. AWS Local Zones are also connected to the parent region via Amazon’s redundant and very high bandwidth private network, giving applications running in AWS Local Zones fast, secure, and seamless access to the rest of AWS services.

### 25. AWS Edge Locations
- **Keyword**: `Content Delivery`, `Low Latency`.
- **Description**: An AWS Edge location is a site that CloudFront uses to cache copies of the content for faster delivery to users at any location.
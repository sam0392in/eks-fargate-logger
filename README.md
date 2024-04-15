# eks-fargate-logger

## Overview

EKS Fargate logger  provides a custom logging solution designed to efficiently ship logs generated by your EKS Fargate applications to a self-hosted Elasticsearch cluster.
By leveraging this solution, you gain centralized log aggregation and improved visibility into the health and performance of your containerized applications running on Amazon EKS with Fargate.

## Key Features

### Simplified Fargate Logging:
Those working with EKS Fargate knows that running a Daemonset on EKS fargate that has access to node for getting container logs is not possible.
This solution eliminates the need to deploy and manage sidecars for log collection and helps to automatically collects the logs centrally.

### Customizable Destination:
Configure your self-hosted Elasticsearch cluster seamlessly as the target for your logs, allowing you to leverage your existing infrastructure.

### Enhanced Visibility:
Aggregate and analyze your container logs in a centralized location, gaining valuable insights into application behavior and potential issues.

## Functioning
![[eks-fargate-logger]]([static/eks-fargate-logger.png])

## Prerequisites
- A running Amazon EKS cluster with Fargate enabled.
- A self-hosted Elasticsearch cluster ready to receive log data.
- Basic familiarity with Kubernetes deployments and configuration management tools like Helm (optional).

## Installation

### AWS Setup
Reference: https://docs.aws.amazon.com/eks/latest/userguide/fargate-logging.html

- Create a dedicated namespace
```
kind: Namespace
apiVersion: v1
metadata:
name: aws-observability
labels:
aws-observability: enabled
```

- Create Fargate profile in EKS, reference: https://docs.aws.amazon.com/eks/latest/userguide/fargate-profile.html

- deploy application

### Deployment Options:

Run on EC2 machine OR as a docker container:
Download the binary from releases and create a docker image / run it directly on server.

#### Run on Kubernetes Cluster:
Helm chart is provided in `k8s` directory with different environments.
NOTE: Create the image and update the helm chart values for image and tag accordingly.

#### Utilise Jenkins Pipeline:
`Jenkinsfile` is also provided in case you use Jenkins for CICD.

- Create a Firehose with HTTP Endpoint. Configure HTTP Endpoint as eks-fargate-logger application endpoint.  reference: https://docs.aws.amazon.com/firehose/latest/dev/create-name.html



## Further Contributions

- I welcome contributions to this project! Feel free to submit pull requests that improve functionality, documentation, or maintainability.
- For questions or feedback on this project, feel free to create an issue in this repository.
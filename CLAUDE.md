# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based Kubernetes demo project with the following components:

- **API Service**: Gin-based HTTP API server (port 8080) that serves as the main entry point
- **Producer Service**: gRPC server (port 50051) that publishes messages to Kafka
- **Consumer Service**: Kafka consumer that processes messages from topics
- **Common Module**: Shared utilities for Kafka and gRPC clients
- **Proto Module**: Protocol buffer definitions for gRPC communication

The services communicate via:
- gRPC for API → Producer communication  
- Kafka for asynchronous message passing between Producer and Consumer
- All services use the shared `common` module for Kafka and gRPC utilities

## Development Commands

### Local Development with Kind
```bash
# Create local Kubernetes cluster
kind create cluster --name demo-cluster --config backend/deployment/demo-cluster.yaml

# Deploy Kafka
kubectl apply -f backend/deployment/kafka-deployment.yaml

# Deploy Kafka UI (access via http://localhost:30080)
kubectl apply -f backend/deployment/kafka-ui.yaml
kubectl port-forward svc/kafka-ui 30080:8080
```

### Building and Deploying Services
```bash
# Build Docker images
docker build -t api:latest -f backend/api/Dockerfile .
docker build -t producer:latest -f backend/kafka/producer/Dockerfile .
docker build -t consumer:latest -f backend/kafka/consumer/Dockerfile .

# Load images into Kind cluster
kind load docker-image api:latest --name demo-cluster
kind load docker-image producer:latest --name demo-cluster  
kind load docker-image consumer:latest --name demo-cluster

# Deploy services
kubectl apply -f backend/deployment/api.yaml
kubectl apply -f backend/deployment/producer.yaml
kubectl apply -f backend/deployment/consumer.yaml
```

### Kafka Operations
```bash
# Connect to Kafka pod
kubectl exec -it $(kubectl get pod -l app=kafka -o jsonpath='{.items[0].metadata.name}') -- bash

# Create topic
kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

# List topics
kafka-topics.sh --list --bootstrap-server localhost:9092
```

### Go Module Structure
The project uses local module replacement in go.mod files:
- `api` depends on `common` and `proto`
- `producer` depends on `common` and `proto` 
- `consumer` depends on `common`
- All modules use `go 1.24.2`

### Terraform (AWS EKS)
```bash
cd terraform/generated/aws/eks
terraform init
terraform plan
terraform apply
```

Deployment order: eks → nat → route_table
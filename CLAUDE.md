# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based Kubernetes demo project. Request flow:

```
Client → (WAF: JP-only) → ALB → API → (gRPC) → Producer → (Kafka) → Consumer
```

- **API Service** (`backend/api`): Gin HTTP server (port 8080). Endpoints: `GET /api/health`, `POST /api/hoge`
- **Producer Service** (`backend/kafka/producer`): gRPC server (port 50051) that publishes to Kafka topic `test-topic`
- **Consumer Service** (`backend/kafka/consumer`): Kafka consumer (group `test-consumer-group`)
- **Trino + OPA**: SQL query engine with OPA-based access control (policies in `backend/deployment/base/opa/`)
- **Common Module** (`backend/common`): `ckafka` (Writer/Reader constructors) and `cgrpc` (connection helpers)
- **Proto Module** (`backend/proto`): protobuf definitions (`user.proto`)

Each Go service follows clean architecture: `internal/domain` / `internal/usecase` / `internal/interface/handler` / `internal/infrastructure`, with `main.go` as the Composition Root (DI via env vars: `KAFKA_BROKER`, `PRODUCER_ADDR`, `USE_SHARED_CONN`).

## Development Commands

### Local Development with Kind
```bash
kind create cluster --name demo-cluster --config backend/deployment/demo-cluster.yaml

# Build images (from repo root)
docker build -t api:latest -f backend/api/Dockerfile .
docker build -t producer:latest -f backend/kafka/producer/Dockerfile .
docker build -t consumer:latest -f backend/kafka/consumer/Dockerfile .

kind load docker-image api:latest producer:latest consumer:latest --name demo-cluster

# Deploy everything (kustomize)
kubectl apply -k backend/deployment/overlays/local

# Verify E2E flow
kubectl port-forward svc/api-service 8080:8080
curl -X POST localhost:8080/api/hoge          # expect 200 + producer hostname
kubectl logs deployment/consumer --tail=5     # expect "received message: ..."
```

### Manifest Structure (kustomize)
```
backend/deployment/
├── base/               # shared definitions for all 6 services
├── overlays/local/     # kind (namespace: default, local images)
├── overlays/develop/   # EKS (namespace: demo, ECR Public, ALB Ingress + WAF)
└── demo-cluster.yaml   # kind cluster config
```
The `${WAF_ACL_ARN}` placeholder in `overlays/develop/ingress.yaml` is substituted by envsubst in CI.

### Kafka Operations
```bash
kubectl exec -it kafka-0 -- bash
/opt/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group test-consumer-group
```
Topic `test-topic` is auto-created. Kafka data is emptyDir (not persisted; intentional for demo).

### Go Module Structure
Local module replacement in go.mod files:
- `api` depends on `common` and `proto`
- `producer` depends on `common` and `proto`
- `consumer` depends on `common` only
- All modules use `go 1.24.2`

Build cache note: in sandboxed shells use `export GOCACHE="$HOME/.cache/go-build"`.

### Terraform (AWS, profile: personal)
```bash
# One-time: tfstate S3 bucket (never destroy)
cd backend/terraform/bootstrap && terraform init && terraform apply

# Per verification session
cd backend/terraform/environments/dev
terraform init && terraform apply    # ~15-20 min
terraform destroy                    # ALWAYS destroy after verification (cost)
```
Creates: VPC (2AZ, single NAT) / EKS `demo-kube-dev` (1.35, t3.large SPOT) / AWS Load Balancer Controller (Helm + Pod Identity) / WAFv2 (JP-only geo allow) / GitHub Actions OIDC role (least privilege) / namespace `demo` (Terraform-managed so the GHA role can be namespace-scoped, and so destroy cascades Ingress→ALB deletion).

ECR Public (`public.ecr.aws/h7f0r1p2/demo-kube`) is intentionally NOT Terraform-managed (images survive destroy).
The GitHub OIDC provider is shared with other apps — referenced via data source, never delete it.

### CI/CD (GitHub Actions)
`develop` push → `.github/workflows/deploy.yml`: paths-filter change detection → matrix build (api/producer/consumer) with GHA cache → ECR Public push → `kubectl apply -k overlays/develop` → rollout → smoke tests (in-cluster curl = 200, consumer log grep, US runner gets 403 from WAF = geo block working).

Required GitHub Variables (set from terraform output): `AWS_DEPLOY_ROLE_ARN`, `EKS_CLUSTER_NAME`, `WAF_ACL_ARN`.

JP-side ALB verification (runner is US-based and gets blocked by WAF by design):
```bash
aws eks update-kubeconfig --name demo-kube-dev --profile personal --region ap-northeast-1
kubectl get ingress api -n demo    # get ALB DNS
curl -X POST http://<ALB_DNS>/api/hoge
```

## Known Gotchas
- Kafka `kafka-storage.sh format` formats the default `log.dirs` (`/tmp/kraft-combined-logs`); the server start override must use the same path
- Consumer must fail-fast (`ckafka.Ping`) if the broker is unreachable at startup — creating a group Reader before the broker is up wedges kafka-go irrecoverably; K8s restart handles retry
- Health probes are intentionally omitted (log noise); the ALB target health check (`/api/health`, 300s interval) cannot be disabled

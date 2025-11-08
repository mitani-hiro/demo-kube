# クラスタ

```sh
kind create cluster --name demo-cluster --config backend/deployment/local/demo-cluster.yaml
```

# kafka デプロイ

```sh
kubectl apply -f backend/deployment/local/kafka-deployment_local.yaml
```

# kafka

```sh
kubectl exec -it $(kubectl get pod -l app=kafka -o jsonpath='{.items[0].metadata.name}') -- bash

# Kafka CLIでトピック作成
kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

# トピック一覧表示
kafka-topics.sh --list --bootstrap-server localhost:9092
```

# kafka ui デプロイ

```sh
kubectl apply -f backend/deployment/local/kafka-ui.yaml
```

# kafka ui 表示

ポートフォワーディングして`http://localhost:30080`にアクセスする。

```sh
kubectl port-forward svc/kafka-ui 30080:8080
```

# swagger ui デプロイ

```sh
kubectl apply -f deployment/swagger-ui.yaml
```

# swagger ui 表示

ポートフォワーディングして`http://localhost:8081`にアクセスする。

```sh
kubectl port-forward svc/swagger-ui 8081:8080
```

# OPA デプロイ

OPA のポリシーとデータを ConfigMap として作成してからデプロイします。

```sh
# ConfigMapを作成（backend/opaディレクトリのファイルから）
kubectl create configmap opa-policies --from-file=trino_access_control.rego=backend/opa/trino_access_control.rego --dry-run=client -o yaml | kubectl apply -f -
kubectl create configmap opa-data --from-file=data.json=backend/opa/data.json --dry-run=client -o yaml | kubectl apply -f -

# OPAをデプロイ
kubectl apply -f backend/deployment/local/opa.yaml

# 確認
kubectl get pods | grep opa
kubectl logs -f deployment/opa
```

ポリシーやデータを更新した場合は、ConfigMap を再作成して OPA ポッドを再起動してください。

```sh
# ConfigMapを更新
kubectl create configmap opa-policies --from-file=trino_access_control.rego=backend/opa/trino_access_control.rego --dry-run=client -o yaml | kubectl apply -f -
kubectl create configmap opa-data --from-file=data.json=backend/opa/data.json --dry-run=client -o yaml | kubectl apply -f -

# OPAポッドを再起動
kubectl rollout restart deployment/opa
```

# ビルド

```sh
docker build -t api:latest -f backend/api/Dockerfile .
docker build -t producer:latest -f backend/kafka/producer/Dockerfile .
docker build -t consumer:latest -f backend/kafka/consumer/Dockerfile .

# 作成済みの場合
kind load docker-image api:latest --name demo-cluster
kubectl delete pod api-server

kind load docker-image producer:latest --name demo-cluster
kubectl delete pod kafka-producer

kind load docker-image consumer:latest --name demo-cluster
kubectl delete pod kafka-consumer

kubectl apply -f backend/deployment/local/api.yaml
kubectl apply -f backend/deployment/local/producer.yaml
kubectl apply -f backend/deployment/local/consumer.yaml

kubectl get pods

```

# terraform

```sh
terraform/generated/aws/eks

terraform init

terraform plan

terraform apply
```

eks -> nat -> root_table

```sh
terraform plan -var="vpc_id=vpc-0f9c1286580033168" -var="nat_gateway_id=nat-05f66b5627bf5b464" -var="subnet_id=subnet-07eb31421fe856a7e"
```

```sh
terraformer import aws \
--resources=eks --regions=ap-northeast-1 --profile=default

terraformer import aws \
--resources=route_table --regions=ap-northeast-1 --profile=default
```

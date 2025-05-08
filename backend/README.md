# クラスタ

```sh
kind create cluster --name demo-cluster --config demo-cluster.yaml
```

# kafka デプロイ

```sh
kubectl apply -f deployment/kafka-deployment.yaml
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
kubectl apply -f deployment/kafka-ui.yaml
```

# kafka ui 表示

ポートフォワーディングして`http://localhost:30080`にアクセスする。

```sh
kubectl port-forward svc/kafka-ui 30080:8080
```

# ビルド

```sh
docker build -t producer:latest -f kafka/producer/Dockerfile .
docker build -t consumer:latest -f kafka/consumer/Dockerfile .

# 作成済みの場合
kind load docker-image producer:latest --name demo-cluster
kubectl delete pod kafka-producer
kind load docker-image consumer:latest --name demo-cluster
kubectl delete pod kafka-consumer

kubectl apply -f deployment/producer.yaml
kubectl apply -f deployment/consumer.yaml

kubectl get pods

```

kind load docker-image producer:latest --name demo-cluster

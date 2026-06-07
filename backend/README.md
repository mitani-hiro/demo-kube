# backend

マイクロサービス構成のK8sデモ。リクエストフロー:

```
クライアント → (WAF: 日本のみ許可) → ALB → API → (gRPC) → Producer → (Kafka) → Consumer
```

| サービス | 役割 | ポート |
|---|---|---|
| api | Gin HTTPサーバー（エントリポイント） | 8080 |
| producer | gRPCサーバー。Kafkaへpublish | 50051 |
| consumer | Kafkaコンシューマー | - |
| kafka | KRaftモードのブローカー | 9092 |
| trino | SQLクエリエンジン（OPAでアクセス制御） | 8080 |
| opa | ポリシー決定エンジン | 8181 |

## ローカル開発（kind）

```sh
# クラスタ作成
kind create cluster --name demo-cluster --config backend/deployment/demo-cluster.yaml

# イメージビルド（リポジトリルートで実行）
docker build -t api:latest -f backend/api/Dockerfile .
docker build -t producer:latest -f backend/kafka/producer/Dockerfile .
docker build -t consumer:latest -f backend/kafka/consumer/Dockerfile .

# kindにロード
kind load docker-image api:latest producer:latest consumer:latest --name demo-cluster

# デプロイ（kustomize）
kubectl apply -k backend/deployment/overlays/local

# 疎通確認
kubectl port-forward svc/api-service 8080:8080
curl -X POST localhost:8080/api/hoge
kubectl logs deployment/consumer --tail=5   # "received message: ..." が出ればE2E疎通OK
```

イメージを更新したら `kind load` 後に `kubectl rollout restart deployment <name>` で反映する。

## マニフェスト構成（kustomize）

```
backend/deployment/
├── base/               # 全サービス共通定義（OPAポリシーは configMapGenerator）
├── overlays/
│   ├── local/          # kind用（namespace: default、ローカルイメージ）
│   └── develop/        # EKS用（namespace: demo、ECR Public、ALB Ingress + WAF）
└── demo-cluster.yaml   # kindクラスタ設定
```

## Kafka操作

```sh
kubectl exec -it kafka-0 -- bash
/opt/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group test-consumer-group
```

トピック `test-topic` は自動作成される（auto.create.topics.enable=true）。

## Trino + OPA

```sh
# ローカル: NodePort 30081 / EKS: port-forward
kubectl port-forward svc/trino 8080:8080
```

OPAポリシー（`backend/deployment/base/opa/`）を更新した場合は再apply（configMapGeneratorがハッシュ付きで再生成）後、`kubectl rollout restart deployment opa trino`。

## Terraform（AWS / EKS）

`AWS_PROFILE=personal` で実行する。

```sh
# 初回のみ: tfstate用S3バケット作成（destroyしない）
cd backend/terraform/bootstrap
terraform init && terraform apply

# 検証開始時: インフラ構築（VPC/EKS/WAF/IAM等、約15-20分）
cd backend/terraform/environments/dev
terraform init
terraform apply

# 検証終了時: コスト削減のため必ず破棄
# ※ Ingress/ALBとnamespaceはTerraform管理のnamespace削除でカスケード削除される
terraform destroy
```

主なリソース: VPC（2AZ・シングルNAT）/ EKS `demo-kube-dev`（t3.large SPOT×1）/ AWS Load Balancer Controller（Helm + Pod Identity）/ WAFv2（Geo: 日本のみ許可）/ GitHub Actions OIDC用IAMロール（最小権限）。

稼働時コスト目安: 約$0.24/h（EKS $0.10 + NAT $0.062 + ノード + ALB + WAF）。

## CI/CD（GitHub Actions）

`develop` へのpushで `.github/workflows/deploy.yml` が実行される:

1. 変更検知（common/proto変更時は全サービス、個別変更時は該当のみ）
2. matrix並列ビルド + GHAキャッシュ → ECR Public push
3. `kubectl apply -k overlays/develop`（WAF ARNはenvsubstで注入）+ rollout
4. 検証: クラスタ内curlで200確認 / consumerログ確認 / USランナーからのALBアクセスが403（WAF動作確認）

必要なGitHub Variables: `AWS_DEPLOY_ROLE_ARN` / `EKS_CLUSTER_NAME` / `WAF_ACL_ARN`（terraform output の値を設定）。

日本からのALB疎通確認:

```sh
aws eks update-kubeconfig --name demo-kube-dev --profile personal --region ap-northeast-1
kubectl get ingress api -n demo   # ALB DNSを取得
curl -X POST http://<ALB_DNS>/api/hoge
```

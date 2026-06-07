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

## クラスタの切り替え（kind ⇔ EKS）

kubeconfig に両クラスタの context を登録して切り替える。

```sh
# EKS の context を登録（terraform apply でクラスタを作り直すたびに再実行する）
AWS_PROFILE=personal aws eks update-kubeconfig --name demo-kube-dev --region ap-northeast-1 --alias demo-kube-dev

# context 一覧と現在の context
kubectl config get-contexts

# 切り替え
kubectl config use-context kind-demo-cluster   # ローカル（kind）
kubectl config use-context demo-kube-dev       # EKS

# 切り替えずに単発で操作する場合は --context を付ける
kubectl --context kind-demo-cluster get pods
kubectl --context demo-kube-dev get pods -n demo
```

※ EKS の認証プロファイル（AWS_PROFILE=personal）は kubeconfig に埋め込まれるため、環境変数の設定は不要。
※ namespace は kind=`default`、EKS=`demo`（EKS 操作時は `-n demo` を付ける）。

## curl での疎通確認

### ローカル（kind）

ALB がないため port-forward 経由で叩く。

```sh
kubectl --context kind-demo-cluster port-forward svc/api-service 8080:8080 &

# ヘルスチェック
curl localhost:8080/api/health
# => {} [200]

# E2E（API → gRPC → Producer → Kafka → Consumer）
curl -X POST localhost:8080/api/hoge
# => {"hostname":"producer-xxxxx","id":999,"name":"Taro Yamada"}

# consumer がメッセージを受信していることを確認
kubectl --context kind-demo-cluster logs deployment/consumer --tail=5
# => "received message: Hello Kafka from Go!" が出ればE2E疎通OK
```

### EKS（ALB / WAF 経由）

ALB の DNS 名は apply のたびに変わるため Ingress から取得する。

```sh
ALB=$(kubectl --context demo-kube-dev get ingress api -n demo \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')

# ヘルスチェック
curl "http://$ALB/api/health"
# => {} [200]

# E2E（WAF → ALB → API → gRPC → Producer → Kafka → Consumer）
curl -X POST "http://$ALB/api/hoge"
# => {"hostname":"producer-xxxxx","id":999,"name":"Taro Yamada"}

# consumer がメッセージを受信していることを確認
kubectl --context demo-kube-dev logs deployment/consumer -n demo --tail=5
# => "received message: Hello Kafka from Go!"
```

※ WAF が日本のみ許可のため、日本国外からのリクエストは 403 になる（GitHub Actions の US ランナーからの 403 はこの仕様を逆手に取った WAF 動作確認）。
※ ALB を経由せず確認したい場合は kind と同様に `kubectl --context demo-kube-dev port-forward svc/api-service 8080:8080 -n demo` でも可。

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
/opt/kafka/bin/kafka-get-offsets.sh --bootstrap-server localhost:9092 --topic test-topic
```

トピック `test-topic` は自動作成される（auto.create.topics.enable=true）。
consumer はグループを使わずパーティション0を直読みする（コーディネーター未準備時のグループ参加スタック回避）。

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

デプロイ後の手動確認は「curl での疎通確認 > EKS」を参照。

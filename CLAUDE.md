# CLAUDE.md

このファイルは Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイドです。

## アーキテクチャ概要

Kubernetes 上のマイクロサービスデモプロジェクト。構成図（アーキテクチャ / Terraform / CI/CD）は `README.md` を参照。

- **api** (`backend/api`): Gin HTTPサーバー。エンドポイントは `GET /api/health` と `POST /api/hoge` のみ
- **producer** (`backend/kafka/producer`): gRPCサーバー。Kafka topic `test-topic` へpublish
- **consumer** (`backend/kafka/consumer`): Kafkaコンシューマー（グループ未使用・パーティション0直読み）
- **trino + opa**: OPAでアクセス制御するSQLクエリエンジン（ポリシーは `backend/deployment/base/opa/`）
- **common** (`backend/common`): `ckafka`（Writer/Readerコンストラクタ）と `cgrpc`（接続ヘルパー）
- **proto** (`backend/proto`): protobuf定義（`user.proto`）

## 開発コマンド

ローカル開発（kind）・クラスタ切り替え・Kafka操作・Terraform運用・CI/CDの具体的なコマンドは `README.md` を、curl疎通手順やTrino+OPAの操作は `backend/README.md` を参照。

開発時の前提:

- 各Goサービスはクリーンアーキテクチャ構成: `internal/domain` / `internal/usecase` / `internal/interface/handler` / `internal/infrastructure`。`main.go` が Composition Root（環境変数 `KAFKA_BROKER`, `PRODUCER_ADDR`, `USE_SHARED_CONN` でDI）
- go.mod の replace でローカルモジュール参照（api/producer は common+proto、consumer は common のみに依存）。全モジュール `go 1.24.2`
- サンドボックスシェルでGoをビルドする場合は `export GOCACHE="$HOME/.cache/go-build"` を設定
- ECR Public（`public.ecr.aws/h7f0r1p2/demo-kube`）は意図的にTerraform管理外（destroyしてもイメージが残る）
- GitHub OIDCプロバイダーは他アプリと共用 — data sourceで参照のみ。**絶対に削除しないこと**
- E2EはCodeBuild（東京）経由でWAFを通過させる（GHAランナーはUS発で403になるため）。テストは `backend/e2e/`（Playwright request API、ブラウザ不要）

## ハマりどころ（Known Gotchas）
- Kafkaの `kafka-storage.sh format` はデフォルトの `log.dirs`（`/tmp/kraft-combined-logs`）をフォーマットする。起動時のoverrideも同じパスを使うこと
- consumerは起動時にbroker不通なら fail-fast（`ckafka.Ping`）すること — broker起動前にグループReaderを作るとkafka-goが復旧不能にスタックする。リトライはK8sの再起動に委ねる（現状はグループ未使用のパーティション直読みでさらに回避）
- liveness/readiness probe は意図的に未設定（ログノイズ対策）。ALBのターゲットヘルスチェック（`/api/health`、間隔300秒）は無効化できない

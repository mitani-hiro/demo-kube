data "aws_caller_identity" "current" {}

# GHA デプロイ用ロールは bootstrap 管理（destroy で消えると CI の AWS 認証が
# すべて失敗するため）。dev からは access entry 用に参照のみ
data "aws_iam_role" "gha_deploy" {
  name = "${local.name_prefix}-gha-deploy-role"
}

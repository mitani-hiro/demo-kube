# GHA デプロイ用 OIDC ロールは destroy しない bootstrap 側で管理する。
# dev 環境（apply/destroy を繰り返す）に置くと destroy 中に GHA の AWS 認証が
# 全て失敗する（ECR push も preflight も不可）ため。
# 参照先リソースの ARN は命名規則から決定的に構成できるので循環依存しない。
data "aws_caller_identity" "current" {}

# 共用 OIDC プロバイダー（他アプリと共用のため参照のみ）
data "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"
}

locals {
  name_prefix = "demo-kube-dev"

  codebuild_e2e_log_group_arn = "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${local.name_prefix}-e2e"
  codebuild_e2e_project_arn   = "arn:aws:codebuild:${var.region}:${data.aws_caller_identity.current.account_id}:project/${local.name_prefix}-e2e"
  eks_cluster_arn             = "arn:aws:eks:${var.region}:${data.aws_caller_identity.current.account_id}:cluster/${local.name_prefix}"
}

module "github_oidc_role" {
  source = "../modules/github_oidc_role"

  codebuild_e2e_log_group_arn = local.codebuild_e2e_log_group_arn
  codebuild_e2e_project_arn   = local.codebuild_e2e_project_arn
  ecr_public_repository_arn   = "arn:aws:ecr-public::${data.aws_caller_identity.current.account_id}:repository/demo-kube"
  eks_cluster_arn             = local.eks_cluster_arn
  github_branch               = "develop"
  github_repository           = "mitani-hiro/demo-kube"
  oidc_provider_arn           = data.aws_iam_openid_connect_provider.github.arn
  role_name                   = "${local.name_prefix}-gha-deploy-role"
}

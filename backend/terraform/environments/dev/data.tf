data "aws_caller_identity" "current" {}

# GitHub Actions OIDC プロバイダーは他アプリと共用のため Terraform 管理外（参照のみ）
data "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"
}

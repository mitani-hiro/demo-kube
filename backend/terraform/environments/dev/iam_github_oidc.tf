module "github_oidc_role" {
  source = "./modules/github_oidc_role"

  codebuild_e2e_log_group_arn = aws_cloudwatch_log_group.e2e.arn
  codebuild_e2e_project_arn   = aws_codebuild_project.e2e.arn
  ecr_public_repository_arn   = var.ecr_public_repository_arn
  eks_cluster_arn             = module.eks.cluster_arn
  github_branch               = var.github_branch
  github_repository           = var.github_repository
  oidc_provider_arn           = data.aws_iam_openid_connect_provider.github.arn
  role_name                   = "${local.name_prefix}-gha-deploy-role"
}

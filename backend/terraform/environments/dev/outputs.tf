output "cluster_name" {
  description = "EKS cluster name (set as GitHub variable EKS_CLUSTER_NAME)"
  value       = module.eks.cluster_name
}

output "gha_deploy_role_arn" {
  description = "GitHub Actions deploy role ARN (set as GitHub variable AWS_DEPLOY_ROLE_ARN)"
  value       = module.github_oidc_role.role_arn
}

output "region" {
  description = "AWS region"
  value       = var.region
}

output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "waf_web_acl_arn" {
  description = "WAFv2 Web ACL ARN to attach to the API ALB (set as GitHub variable WAF_ACL_ARN)"
  value       = aws_wafv2_web_acl.api.arn
}

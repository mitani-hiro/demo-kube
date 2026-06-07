output "cluster_name" {
  description = "EKS cluster name (set as GitHub variable EKS_CLUSTER_NAME)"
  value       = module.eks.cluster_name
}

output "e2e_codebuild_project_name" {
  description = "CodeBuild project for Playwright E2E (set as GitHub variable E2E_CODEBUILD_PROJECT)"
  value       = aws_codebuild_project.e2e.name
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

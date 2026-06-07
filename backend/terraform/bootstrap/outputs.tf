output "gha_deploy_role_arn" {
  description = "GitHub Actions deploy role ARN (set as GitHub variable AWS_DEPLOY_ROLE_ARN)"
  value       = module.github_oidc_role.role_arn
}

output "state_bucket_name" {
  description = "Name of the S3 bucket that stores Terraform state"
  value       = aws_s3_bucket.tfstate.id
}

variable "codebuild_e2e_log_group_arn" {
  description = "ARN of the CloudWatch log group of the E2E CodeBuild project (for streaming build logs)"
  type        = string
}

variable "codebuild_e2e_project_arn" {
  description = "ARN of the E2E CodeBuild project the role may start"
  type        = string
}

variable "ecr_public_repository_arn" {
  description = "ARN of the ECR Public repository the role may push images to"
  type        = string
}

variable "eks_cluster_arn" {
  description = "ARN of the EKS cluster the role may describe (for update-kubeconfig)"
  type        = string
}

variable "github_branch" {
  description = "Branch allowed to assume this role"
  type        = string
}

variable "github_repository" {
  description = "GitHub repository (owner/name) allowed to assume this role"
  type        = string
}

variable "oidc_provider_arn" {
  description = "ARN of the shared GitHub Actions OIDC provider"
  type        = string
}

variable "role_name" {
  description = "Name of the IAM role"
  type        = string
}

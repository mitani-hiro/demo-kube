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

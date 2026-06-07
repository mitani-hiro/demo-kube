variable "app_namespace" {
  description = "Kubernetes namespace for application workloads"
  type        = string
  default     = "demo"
}

variable "ecr_public_repository_arn" {
  description = "ARN of the ECR Public repository used as the image registry (managed outside Terraform so images survive destroy)"
  type        = string
  default     = "arn:aws:ecr-public::421156463971:repository/demo-kube"
}

variable "github_branch" {
  description = "Branch allowed to assume the GitHub Actions deploy role"
  type        = string
  default     = "develop"
}

variable "github_repository" {
  description = "GitHub repository (owner/name) allowed to assume the deploy role"
  type        = string
  default     = "mitani-hiro/demo-kube"
}

variable "kubernetes_version" {
  description = "EKS Kubernetes version (must be within standard support)"
  type        = string
  default     = "1.35"
}

variable "profile" {
  description = "AWS CLI profile name"
  type        = string
  default     = "personal"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

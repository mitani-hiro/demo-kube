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

variable "state_bucket_name" {
  description = "Name of the S3 bucket that stores Terraform state"
  type        = string
  default     = "demo-kube-tfstate-421156463971-apne1"
}

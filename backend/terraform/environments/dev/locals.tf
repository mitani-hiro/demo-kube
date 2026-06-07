locals {
  azs         = ["${var.region}a", "${var.region}c"]
  name_prefix = "demo-kube-dev"
  vpc_cidr    = "10.0.0.0/16"

  common_tags = {
    Project     = "demo-kube"
    Environment = "dev"
    ManagedBy   = "terraform"
    Repository  = "github.com/${var.github_repository}"
  }
}

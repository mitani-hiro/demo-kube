provider "aws" {
  region  = var.region
  profile = var.profile

  default_tags {
    tags = {
      Project     = "demo-kube"
      Environment = "shared"
      ManagedBy   = "terraform"
      Repository  = "github.com/mitani-hiro/demo-kube"
    }
  }
}

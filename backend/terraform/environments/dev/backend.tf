terraform {
  backend "s3" {
    bucket       = "demo-kube-tfstate-421156463971-apne1"
    key          = "dev/terraform.tfstate"
    region       = "ap-northeast-1"
    profile      = "personal"
    encrypt      = true
    use_lockfile = true # S3 ネイティブロック（DynamoDB ロックは非推奨）
  }
}

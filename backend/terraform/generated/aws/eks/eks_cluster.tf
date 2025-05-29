resource "aws_eks_cluster" "tfer--eks-cluster-demo" {
  kubernetes_network_config {
    ip_family         = "ipv4"
    service_ipv4_cidr = "172.20.0.0/16"
  }

  name     = "eks-cluster-demo"
  role_arn = "arn:aws:iam::421156463971:role/EKS-ClusterRole-demo"
  version  = "1.32"

  vpc_config {
    endpoint_private_access = "true"
    endpoint_public_access  = "true"
    public_access_cidrs     = ["0.0.0.0/0"]
    security_group_ids      = ["sg-0f6e6f1182888337e"]
    subnet_ids              = ["subnet-0437fba6aa7d6fee9", "subnet-07d60709f844e57e3"]
  }
}

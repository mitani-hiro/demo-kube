# AWS Load Balancer Controller の権限は EKS Pod Identity で付与（IRSA より新しい推奨方式）
module "lb_controller_pod_identity" {
  source  = "terraform-aws-modules/eks-pod-identity/aws"
  version = "~> 2.0"

  name = "${local.name_prefix}-lbc"

  attach_aws_lb_controller_policy = true

  associations = {
    lbc = {
      cluster_name    = module.eks.cluster_name
      namespace       = "kube-system"
      service_account = "aws-load-balancer-controller"
    }
  }
}

resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "3.4.0"
  namespace  = "kube-system"

  # ノードの IMDS hop limit が 1 のため Pod から取得できない region / vpcId は明示する
  set = [
    {
      name  = "clusterName"
      value = module.eks.cluster_name
    },
    {
      name  = "region"
      value = var.region
    },
    {
      name  = "vpcId"
      value = module.vpc.vpc_id
    },
    {
      name  = "replicaCount"
      value = "1"
    },
    {
      name  = "serviceAccount.name"
      value = "aws-load-balancer-controller"
    },
  ]

  depends_on = [
    module.eks,
    module.lb_controller_pod_identity,
  ]
}

# アプリ namespace は Terraform 管理にする:
# - GHA ロールの権限を namespace スコープに限定できる（namespace 作成権限が不要になる）
# - destroy 時に namespace 削除 → Ingress/ALB のカスケード削除が LBC 稼働中に行われる
#   （depends_on により LBC より先に削除されるため ALB が孤児にならない）
resource "kubernetes_namespace" "app" {
  metadata {
    name = var.app_namespace
  }

  depends_on = [helm_release.aws_load_balancer_controller]
}

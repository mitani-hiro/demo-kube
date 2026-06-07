# API 用 WAF: 日本以外からのアクセスをブロック（Geo 制限のみ）
# ALB への関連付けは Ingress アノテーション alb.ingress.kubernetes.io/wafv2-acl-arn で行う
resource "aws_wafv2_web_acl" "api" {
  name        = "${local.name_prefix}-api-waf"
  description = "Allow requests from Japan only"
  scope       = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "allow-jp"
    priority = 1

    action {
      allow {}
    }

    statement {
      geo_match_statement {
        country_codes = ["JP"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${local.name_prefix}-allow-jp"
      sampled_requests_enabled   = true
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "${local.name_prefix}-api-waf"
    sampled_requests_enabled   = true
  }
}

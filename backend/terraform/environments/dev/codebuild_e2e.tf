# E2E テスト用 CodeBuild プロジェクト。
# GitHub Actions ホストランナー（US）は WAF の Geo 制限（日本のみ許可）で
# ブロックされるため、東京リージョンの CodeBuild から Playwright を実行して
# WAF を通過した E2E テストを成立させる。
resource "aws_cloudwatch_log_group" "e2e" {
  name              = "/aws/codebuild/${local.name_prefix}-e2e"
  retention_in_days = 14
}

data "aws_iam_policy_document" "codebuild_assume" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["codebuild.amazonaws.com"]
    }

    # confused deputy 対策
    condition {
      test     = "StringEquals"
      variable = "aws:SourceAccount"
      values   = [data.aws_caller_identity.current.account_id]
    }
  }
}

resource "aws_iam_role" "codebuild_e2e" {
  name               = "${local.name_prefix}-e2e-codebuild-role"
  description        = "CodeBuild service role for Playwright E2E (logs only)"
  assume_role_policy = data.aws_iam_policy_document.codebuild_assume.json
}

data "aws_iam_policy_document" "codebuild_e2e" {
  statement {
    sid = "CloudWatchLogs"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["${aws_cloudwatch_log_group.e2e.arn}:*"]
  }
}

resource "aws_iam_role_policy" "codebuild_e2e" {
  name   = "logs"
  role   = aws_iam_role.codebuild_e2e.id
  policy = data.aws_iam_policy_document.codebuild_e2e.json
}

resource "aws_codebuild_project" "e2e" {
  name          = "${local.name_prefix}-e2e"
  description   = "Playwright API E2E from Tokyo region (passes WAF geo JP restriction)"
  service_role  = aws_iam_role.codebuild_e2e.arn
  build_timeout = 10

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    type         = "LINUX_CONTAINER"
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/standard:7.0"
  }

  # public リポジトリのため認証不要
  source {
    type            = "GITHUB"
    location        = "https://github.com/${var.github_repository}.git"
    buildspec       = "backend/e2e/buildspec.yml"
    git_clone_depth = 1
  }

  # start-build の --source-version で上書きされる既定値
  source_version = "refs/heads/${var.github_branch}"

  logs_config {
    cloudwatch_logs {
      group_name = aws_cloudwatch_log_group.e2e.name
    }
  }
}

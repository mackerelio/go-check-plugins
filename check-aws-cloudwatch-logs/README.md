# check-aws-cloudwatch-logs

## Description

Checks Amazon CloudWatch Logs using Filter Syntax.

https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html

## Setting

```toml
[plugin.checks.aws-cloudwatch-logs-sample]
command = """
/path/to/check-aws-cloudwatch-logs \
  --log-group-name /aws/lambda/sample_log_group \
  --pattern "Error" \
  --critical-over M \
  --warning-over N \
"""
env = { AWS_REGION = "ap-northeast-1" }
```

Note that for `--pattern` argument, you can use the syntax described in [Filter and Pattern Syntax - Amazon CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html). This is not a regular expression.

The plugin uses the instance profile if possible, or you can configure `AWS_PROFILE` or `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables in the `env` settings.

## Required action
Following action on the target log group is required to perform the monitoring.

- `logs:FilterLogEvents`

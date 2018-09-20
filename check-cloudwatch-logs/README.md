# check-cloudwatch-logs

## Description

Checks CloudWatch Logs using Filter Syntax.

https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html

## Setting

```toml
[plugin.checks.cloudwatch-logs-sample]
command = """
/path/to/check-cloudwatch-logs \
  --log-group-name /aws/lambda/sample_log_group \
  --pattern "Error" \
  --critical-over M \
  --warning-over N \
"""
env = { AWS_REGION = "ap-northeast-1" }
```

Note that for `--pattern` argument, you can use the syntax described in [Filter and Pattern Syntax - Amazon CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html). This is not a regular expression.

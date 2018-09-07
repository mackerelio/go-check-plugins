# check-cloudwatch-logs

## Description

Checks CloudWatch Logs using Filter Syntax.

https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html

## Setting

```toml
[plugin.checks.cloudwatch-logs-sample]
command = """
/path/to/check-cloudwatch-logs \
  # --access-key-id ... --secret-access-key ...  (not necessary when using the instance profile) \
  --region ap-northeast-1 \
  --log-group-name /aws/lambda/sample_log_group \
  --pattern "Error" \
  --critical-over M \
  --warning-over N \
"""
```

Note that for `--pattern` argument, you can use the syntax described in [Filter and Pattern Syntax - Amazon CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html). This is not a regular expression.

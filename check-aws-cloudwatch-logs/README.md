# check-aws-cloudwatch-logs

## Description
Checks Amazon CloudWatch Logs using Filter Syntax.

https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html

## Synopsis
```
check-aws-cloudwatch-logs --log-group-name /aws/lambda/sample_log_group --pattern "Error" --critical-over 10 --warning-over 5
```

## Required action
Following action on the target log group is required to perform the monitoring.

- `logs:FilterLogEvents`

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-aws-cloudwatch-logs
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-aws-cloudwatch-logs --log-group-name /aws/lambda/sample_log_group --pattern "Error" --critical-over 10 --warning-over 5
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.aws-cloudwatch-logs-sample]
command = ["check-aws-cloudwatch-logs", "--log-group-name", "/aws/lambda/sample_log_group", "--pattern", "Error", "--critical-over", "10", "--warning-over", "5"]
```

## Usage
### Options

```
      --log-group-name=LOG-GROUP-NAME                    Log group name
      --log-stream-name-prefix=LOG-STREAM-NAME-PREFIX    Log stream name prefix
  -p, --pattern=PATTERN                                  Pattern to search for. The value is recognized as the pattern syntax of CloudWatch Logs.
  -w, --warning-over=WARNING                             Trigger a warning if matched lines is over a number
  -c, --critical-over=CRITICAL                           Trigger a critical if matched lines is over a number
  -s, --state-dir=DIR                                    Dir to keep state files under
  -r, --return                                           Output matched lines
  -t, --max-retries=MAX-RETRIES                          Maximum number of retries to call the AWS API
```

Note that for `--pattern` argument, you can use the syntax described in [Filter and Pattern Syntax - Amazon CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html). This is not a regular expression.

The plugin uses the instance profile if possible, or you can configure `AWS_PROFILE` or `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables in the `env` settings.

## For more information
Please refer to the following.

- Read Mackerel Docs; [Monitoring Amazon CloudWatch Logs - Mackerel Docs](https://mackerel.io/docs/entry/howto/check/aws-cloudwatch-logs)
- Execute `check-aws-cloudwatch-logs -h` and you can get command line options.

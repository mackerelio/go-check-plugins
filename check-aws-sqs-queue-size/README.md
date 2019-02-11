# check-aws-sqs-queue-size

## Description
Check the queue size of specified SQS queue.

## Synopsis
```
check-aws-sqs-queue-size --access-key-id=<aws-access-key-id> --secret-access-key=<aws-secret-access-key> --region=<aws-region> --queue=<the-name-of-the-aws-sqs-queue>
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-aws-sqs-queue-size
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-aws-sqs-queue-size --access-key-id=<aws-access-key-id> --secret-access-key=<aws-secret-access-key> --region=<aws-region> --queue=<the-name-of-the-aws-sqs-queue>
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.aws-sqs-queue-size-sample]
command = ["check-aws-sqs-queue-size", "--access-key-id", "<aws-access-key-id>", "--secret-access-key", "<aws-secret-access-key>", "--region", "<aws-region>", "--queue", "<the-name-of-the-aws-sqs-queue>"]
```

## Usage
### Options

```
  -r, --region=            AWS Region
  -i, --access-key-id=     AWS Access Key ID
  -s, --secret-access-key= AWS Secret Access Key
  -q, --queue=             The name of the queue name
  -w, --warning=           warning if the number of queues is over (default: 10)
  -c, --critical=          critical if the number of queues is over (default: 100)
```

## For more information
Please execute `check-aws-sqs-queue-size -h` and you can get command line options.

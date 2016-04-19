# check-aws-sqs-queue-size

## Description

Check the queue size of specified SQS queue.

## Setting

```
[plugin.checks.aws-sqs-queue-size]
command = "/path/to/check-aws-sqs-queue-size --access-key-id=<aws-access-key-id> --secret-access-key=<aws-secret-access-key> --region=<aws-region> --queue=<the-name-of-the-aws-sqs-queue>"
```

## Options

```
-r --region            The AWS region (required)
-i --access-key-id     The AWS Access Key ID (required)
-s --secret-access-key The AWS Secret Access Key (required)
-q --queue             The name of the queue (required)
-w --warning           The warning threshold of the queue size (default: 10)
-c --critical          The critical threshold of the queue size (default: 100)
```

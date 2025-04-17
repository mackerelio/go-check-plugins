package checkawssqsqueuesize

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "SQSQueueSize"
	ckr.Exit()
}

var opts struct {
	Region          string `short:"r" long:"region" description:"AWS Region"`
	AccessKeyID     string `short:"i" long:"access-key-id" description:"AWS Access Key ID"`
	SecretAccessKey string `short:"s" long:"secret-access-key" description:"AWS Secret Access Key"`
	QueueName       string `short:"q" long:"queue" required:"true" description:"The name of the queue name"`
	Warn            int    `short:"w" long:"warning" default:"10" description:"warning if the number of queues is over"`
	Crit            int    `short:"c" long:"critical" default:"100" description:"critical if the number of queues is over"`
}

const sqsAttributeOfQueueSize = "ApproximateNumberOfMessages"

func createService(ctx context.Context, region, awsAccessKeyID, awsSecretAccessKey string) (*sqs.Client, error) {
	var opts []func(*config.LoadOptions) error
	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")))
	}
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}

func getSqsQueueSize(ctx context.Context, region, awsAccessKeyID, awsSecretAccessKey, queueName string) (int, error) {
	sqsClient, err := createService(ctx, region, awsAccessKeyID, awsSecretAccessKey)
	if err != nil {
		return -1, err
	}

	// Get queue url
	q, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return -1, err
	}

	// Get queue attribute
	attr, err := sqsClient.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: q.QueueUrl,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
		},
	})
	if err != nil {
		return -1, err
	}

	// Queue size
	sizeStr, ok := attr.Attributes[sqsAttributeOfQueueSize]
	if !ok {
		return -1, errors.New("attribute not found")
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return -1, err
	}

	return size, nil
}

func run(args []string) *checkers.Checker {
	ctx := context.Background()

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	size, err := getSqsQueueSize(ctx, opts.Region, opts.AccessKeyID, opts.SecretAccessKey, opts.QueueName)
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, err.Error())
	}

	var chkSt checkers.Status
	var msg string
	if opts.Crit < size {
		msg = fmt.Sprintf("size %d > %d in %s", size, opts.Crit, opts.QueueName)
		chkSt = checkers.CRITICAL
	} else if opts.Warn < size {
		msg = fmt.Sprintf("size %d > %d in %s", size, opts.Warn, opts.QueueName)
		chkSt = checkers.WARNING
	} else {
		msg = fmt.Sprintf("size %d < warning %d, critical %d in %s", size, opts.Warn, opts.Crit, opts.QueueName)
		chkSt = checkers.OK
	}

	return checkers.NewChecker(chkSt, msg)
}

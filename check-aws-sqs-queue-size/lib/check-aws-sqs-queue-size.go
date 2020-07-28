package checkawssqsqueuesize

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func createService(region, awsAccessKeyID, awsSecretAccessKey string) (*sqs.SQS, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig()
	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""))
	}
	if region != "" {
		config = config.WithRegion(region)
	}
	return sqs.New(sess, config), nil
}

func getSqsQueueSize(region, awsAccessKeyID, awsSecretAccessKey, queueName string) (int, error) {
	sqsClient, err := createService(region, awsAccessKeyID, awsSecretAccessKey)
	if err != nil {
		return -1, err
	}

	// Get queue url
	q, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return -1, err
	}

	// Get queue attribute
	attr, err := sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl:       q.QueueUrl,
		AttributeNames: []*string{aws.String(sqsAttributeOfQueueSize)},
	})
	if err != nil {
		return -1, err
	}

	// Queue size
	sizeStr, ok := attr.Attributes[sqsAttributeOfQueueSize]
	if !ok || sizeStr == nil {
		return -1, errors.New("attribute not found")
	}
	size, err := strconv.Atoi(*sizeStr)
	if err != nil {
		return -1, err
	}

	return size, nil
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	size, err := getSqsQueueSize(opts.Region, opts.AccessKeyID, opts.SecretAccessKey, opts.QueueName)
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

package checkawssqsqueuesize

import (
	"fmt"
	"os"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/sqs"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"

	"strconv"
	"time"
)

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "SQSQueueSize"
	ckr.Exit()
}

var opts struct {
	Region          string `short:"r" long:"region" required:"true" description:"AWS Region"`
	AccessKeyID     string `short:"i" long:"access-key-id" required:"true" description:"AWS Access Key ID"`
	SecretAccessKey string `short:"s" long:"secret-access-key" required:"true" description:"AWS Secret Access Key"`
	QueueName       string `short:"q" long:"queue" required:"true" description:"The name of the queue name"`
	Warn            int    `short:"w" long:"warning" default:"10" description:"warning if the number of queues is over"`
	Crit            int    `short:"c" long:"critical" default:"100" description:"critical if the number of queues is over"`
}

const sqsAttributeOfQueueSize = "ApproximateNumberOfMessages"

func getSqsQueueSize(region, awsAccessKeyID, awsSecretAccessKey, queueName string) (int, error) {
	// Auth
	auth, err := aws.GetAuth(awsAccessKeyID, awsSecretAccessKey, "", time.Now())
	if err != nil {
		return -1, err
	}

	// SQS
	sqsClient := sqs.New(auth, aws.GetRegion(region))
	queue, err := sqsClient.GetQueue(queueName)
	if err != nil {
		return -1, err
	}

	// Get queue attribute
	attr, err := queue.GetQueueAttributes(sqsAttributeOfQueueSize)
	if err != nil {
		return -1, err
	}

	// Queue size
	size, err := strconv.Atoi(attr.Attributes[0].Value)
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

package checkcloudwatchlogs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	"github.com/mackerelio/checkers"
)

type logOpts struct {
	Region          string `long:"region" value-name:"REGION" description:"AWS Region"`
	AccessKeyID     string `long:"access-key-id" value-name:"ACCESS-KEY-ID" description:"AWS Access Key ID"`
	SecretAccessKey string `long:"secret-access-key" value-name:"SECRET-ACCESS-KEY" description:"AWS Secret Access Key"`
	LogGroupName    string `long:"log-group-name" value-name:"LOG-GROUP-NAME" description:"Log group name"`

	Pattern string `long:"pattern" required:"true" value-name:"PATTERN" description:"Pattern to search for. The value is recognized as the pattern syntax of CloudWatch Logs."`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "CloudWatch Logs"
	ckr.Exit()
}

type cloudwatchLogsPlugin struct {
	Service      cloudwatchlogsiface.CloudWatchLogsAPI
	LogGroupName string
	Pattern      string
}

func newCloudwatchLogsPlugin(args []string) (*cloudwatchLogsPlugin, error) {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		return nil, err
	}
	if opts.LogGroupName == "" {
		return nil, errors.New("specify log group name")
	}
	service, err := createService(opts)
	if err != nil {
		return nil, err
	}
	return &cloudwatchLogsPlugin{
		Service:      service,
		LogGroupName: opts.LogGroupName,
		Pattern:      opts.Pattern,
	}, nil
}

func createService(opts *logOpts) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	config := aws.NewConfig()
	if opts.AccessKeyID != "" && opts.SecretAccessKey != "" {
		config = config.WithCredentials(
			credentials.NewStaticCredentials(opts.AccessKeyID, opts.SecretAccessKey, ""),
		)
	}
	if opts.Region != "" {
		config = config.WithRegion(opts.Region)
	}
	return cloudwatchlogs.New(sess, config), nil
}

func (p *cloudwatchLogsPlugin) run() ([]string, error) {
	var nextToken *string
	var messages []string
	for {
		startTime := time.Now().Add(-5 * time.Minute)
		output, err := p.Service.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
			StartTime:     aws.Int64(startTime.Unix() * 1000),
			LogGroupName:  aws.String(p.LogGroupName),
			NextToken:     nextToken,
			FilterPattern: aws.String(p.Pattern),
		})
		if err != nil {
			return nil, err
		}
		for _, ev := range output.Events {
			messages = append(messages, *ev.Message)
		}
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
		time.Sleep(250 * time.Millisecond)
	}
	return messages, nil
}

func run(args []string) *checkers.Checker {
	p, err := newCloudwatchLogsPlugin(args)
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, fmt.Sprint(err))
	}
	messages, err := p.run()
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, fmt.Sprint(err))
	}
	if messages != nil {
		return checkers.NewChecker(checkers.WARNING, strings.Join(messages, ""))
	}
	return checkers.NewChecker(checkers.OK, "ok")
}

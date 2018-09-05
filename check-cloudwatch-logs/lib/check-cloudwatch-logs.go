package checkcloudwatchlogs

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	"github.com/mackerelio/checkers"
)

type logOpts struct {
	Region          string `long:"region" value-name:"REGION" description:"AWS Region"`
	AccessKeyID     string `long:"access-key-id" value-name:"ACCESS-KEY-ID" description:"AWS Access Key ID"`
	SecretAccessKey string `long:"secret-access-key" value-name:"SECRET-ACCESS-KEY" description:"AWS Secret Access Key"`
	LogGroupName    string `long:"log-group-name" value-name:"LOG-GROUP-NAME" description:"Log group name"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "CloudWatch Logs"
	ckr.Exit()
}

type cloudwatchLogsPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	LogGroupName    string
}

func newCloudwatchLogsPlugin(args []string) (*cloudwatchLogsPlugin, error) {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		return nil, err
	}
	return &cloudwatchLogsPlugin{
		Region:          opts.Region,
		AccessKeyID:     opts.AccessKeyID,
		SecretAccessKey: opts.SecretAccessKey,
		LogGroupName:    opts.LogGroupName,
	}, nil
}

func (p *cloudwatchLogsPlugin) getService() (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(
			credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""),
		)
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}
	return cloudwatchlogs.New(sess, config), nil
}

func (p *cloudwatchLogsPlugin) run() error {
	if p.LogGroupName == "" {
		return errors.New("specify log group name")
	}
	service, err := p.getService()
	if err != nil {
		return err
	}
	now := time.Now().Add(-3 * time.Minute)
	events, err := service.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
		StartTime:    aws.Int64(now.UnixNano()),
		LogGroupName: aws.String(p.LogGroupName),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", events)
	return nil
}

func run(args []string) *checkers.Checker {
	p, err := newCloudwatchLogsPlugin(args)
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, fmt.Sprint(err))
	}
	err = p.run()
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, fmt.Sprint(err))
	}
	return checkers.NewChecker(checkers.OK, "ok")
}

package checkcloudwatchlogs

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type logOpts struct {
	Region          string `long:"region" value-name:"REGION" description:"AWS Region"`
	AccessKeyID     string `long:"access-key-id" value-name:"ACCESS-KEY-ID" description:"AWS Access Key ID"`
	SecretAccessKey string `long:"secret-access-key" value-name:"SECRET-ACCESS-KEY" description:"AWS Secret Access Key"`
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
}

func newCloudwatchLogsPlugin(opts *logOpts) *cloudwatchLogsPlugin {
	return &cloudwatchLogsPlugin{
		Region:          opts.Region,
		AccessKeyID:     opts.AccessKeyID,
		SecretAccessKey: opts.SecretAccessKey,
	}
}

func (p *cloudwatchLogsPlugin) getClient() (*cloudwatchlogs.CloudWatchLogs, error) {
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

func run(args []string) *checkers.Checker {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		fmt.Printf("%#v\n", err)
		os.Exit(1)
	}
	p := newCloudwatchLogsPlugin(opts)
	s, _ := p.getClient()
	fmt.Printf("%#v\n", s)
	return checkers.NewChecker(checkers.OK, "ok")
}

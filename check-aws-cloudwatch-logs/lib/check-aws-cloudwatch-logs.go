package checkawscloudwatchlogs

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/jessevdk/go-flags"

	"github.com/mackerelio/checkers"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/natefinch/atomic"
)

// overwritten with syscall.SIGTERM on unix environment (see check-log_unix.go)
var defaultSignal = os.Interrupt

type logOpts struct {
	LogGroupName        string `long:"log-group-name" required:"true" value-name:"LOG-GROUP-NAME" description:"Log group name" unquote:"false"`
	LogStreamNamePrefix string `long:"log-stream-name-prefix" value-name:"LOG-STREAM-NAME-PREFIX" description:"Log stream name prefix" unquote:"false"`

	Pattern       string `short:"p" long:"pattern" required:"true" value-name:"PATTERN" description:"Pattern to search for. The value is recognized as the pattern syntax of CloudWatch Logs." unquote:"false"`
	WarningOver   int    `short:"w" long:"warning-over" value-name:"WARNING" description:"Trigger a warning if matched lines is over a number"`
	CriticalOver  int    `short:"c" long:"critical-over" value-name:"CRITICAL" description:"Trigger a critical if matched lines is over a number"`
	StateDir      string `short:"s" long:"state-dir" value-name:"DIR" description:"Dir to keep state files under" unquote:"false"`
	ReturnContent bool   `short:"r" long:"return" description:"Output matched lines"`
	MaxRetries    int    `short:"t" long:"max-retries" value-name:"MAX-RETRIES" description:"Maximum number of retries to call the AWS API"`
}

// Do the plugin
func Do() {
	ctx, stop := signal.NotifyContext(context.Background(), defaultSignal)
	defer stop()

	ckr := run(ctx, os.Args[1:])
	ckr.Name = "CloudWatch Logs"
	ckr.Exit()
}

type awsCloudwatchLogsPlugin struct {
	Service   cloudwatchlogsiface.CloudWatchLogsAPI
	StateFile string
	*logOpts
}

func newCloudwatchLogsPlugin(opts *logOpts, args []string) (*awsCloudwatchLogsPlugin, error) {
	var err error
	p := &awsCloudwatchLogsPlugin{logOpts: opts}
	p.Service, err = createService(opts)
	if err != nil {
		return nil, err
	}
	if p.StateDir == "" {
		workdir := pluginutil.PluginWorkDir()
		p.StateDir = filepath.Join(workdir, "check-cloudwatch-logs")
	}
	p.StateFile = getStateFile(p.StateDir, opts.LogGroupName, opts.LogStreamNamePrefix, args)
	return p, nil
}

var stateRe = regexp.MustCompile(`[^-a-zA-Z0-9_.]`)

func getStateFile(stateDir, logGroupName, logStreamNamePrefix string, args []string) string {
	return filepath.Join(
		stateDir,
		fmt.Sprintf(
			"%s-%x.json",
			strings.TrimLeft(stateRe.ReplaceAllString(logGroupName+"_"+logStreamNamePrefix, "_"), "_"),
			md5.Sum([]byte(
				strings.Join(
					[]string{
						os.Getenv("AWS_PROFILE"),
						os.Getenv("AWS_ACCESS_KEY_ID"),
						os.Getenv("AWS_SECRET_ACCESS_KEY"),
						os.Getenv("AWS_REGION"),
						strings.Join(args, " "),
					},
					" ",
				)),
			),
		),
	)
}

func createAWSConfig(opts *logOpts) *aws.Config {
	conf := aws.NewConfig()
	if opts.MaxRetries > 0 {
		return conf.WithMaxRetries(opts.MaxRetries)
	}
	return conf
}

func createService(opts *logOpts) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.New(sess, createAWSConfig(opts)), nil
}

type logState struct {
	NextToken *string
	StartTime *int64
}

func (p *awsCloudwatchLogsPlugin) collect(ctx context.Context, now time.Time) ([]string, error) {
	s, err := p.loadState()
	if err != nil {
		return nil, err
	}
	if s.StartTime == nil || *s.StartTime <= now.Add(-time.Hour).Unix()*1000 {
		s.StartTime = aws.Int64(now.Add(-1*time.Minute).Unix() * 1000)
		s.NextToken = nil
	}
	var messages []string
	input := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:     s.StartTime,
		LogGroupName:  aws.String(p.LogGroupName),
		NextToken:     s.NextToken,
		FilterPattern: aws.String(p.Pattern),
	}
	if p.LogStreamNamePrefix != "" {
		input.LogStreamNamePrefix = aws.String(p.LogStreamNamePrefix)
	}
	err = p.Service.FilterLogEventsPages(input, func(output *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		if ctx.Err() != nil {
			return false
		}
		for _, event := range output.Events {
			messages = append(messages, *event.Message)
			if s.StartTime == nil || *s.StartTime <= *event.Timestamp {
				s.StartTime = aws.Int64(*event.Timestamp + 1)
			}
		}
		s.NextToken = output.NextToken
		if lastPage {
			s.NextToken = nil
		}
		err = p.saveState(s)
		if err != nil {
			return false
		}
		time.Sleep(150 * time.Millisecond)
		return true
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (p *awsCloudwatchLogsPlugin) loadState() (*logState, error) {
	var s logState
	f, err := os.Open(p.StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &s, nil
		}
		return nil, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (p *awsCloudwatchLogsPlugin) saveState(s *logState) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(s); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p.StateFile), 0755); err != nil {
		return err
	}
	return atomic.WriteFile(p.StateFile, &buf)
}

func (p *awsCloudwatchLogsPlugin) check(messages []string) *checkers.Checker {
	status := checkers.OK
	msg := fmt.Sprint(len(messages))
	if len(messages) > p.CriticalOver {
		status = checkers.CRITICAL
		msg += " > " + fmt.Sprint(p.CriticalOver)
	} else if len(messages) > p.WarningOver {
		status = checkers.WARNING
		msg += " > " + fmt.Sprint(p.WarningOver)
	}
	msg += " messages for pattern /" + p.Pattern + "/"
	if status != checkers.OK && p.ReturnContent {
		msg += "\n" + strings.Join(messages, "")
	}
	return checkers.NewChecker(status, msg)
}

func (p *awsCloudwatchLogsPlugin) run(ctx context.Context, now time.Time) *checkers.Checker {
	messages, err := p.collect(ctx, now)
	if err != nil {
		return checkers.Unknown(fmt.Sprint(err))
	}
	return p.check(messages)
}

func run(ctx context.Context, args []string) *checkers.Checker {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		os.Exit(1)
	}
	p, err := newCloudwatchLogsPlugin(opts, args)
	if err != nil {
		return checkers.Unknown(fmt.Sprint(err))
	}
	return p.run(ctx, time.Now())
}

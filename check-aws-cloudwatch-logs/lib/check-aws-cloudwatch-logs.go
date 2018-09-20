package checkawscloudwatchlogs

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
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
)

type logOpts struct {
	LogGroupName        string `long:"log-group-name" required:"true" value-name:"LOG-GROUP-NAME" description:"Log group name"`
	LogStreamNamePrefix string `long:"log-stream-name-prefix" value-name:"LOG-STREAM-NAME-PREFIX" description:"Log stream name prefix"`

	Pattern       string `short:"p" long:"pattern" required:"true" value-name:"PATTERN" description:"Pattern to search for. The value is recognized as the pattern syntax of CloudWatch Logs."`
	WarningOver   int    `short:"w" long:"warning-over" value-name:"WARNING" description:"Trigger a warning if matched lines is over a number"`
	CriticalOver  int    `short:"c" long:"critical-over" value-name:"CRITICAL" description:"Trigger a critical if matched lines is over a number"`
	StateDir      string `short:"s" long:"state-dir" value-name:"DIR" description:"Dir to keep state files under"`
	ReturnContent bool   `short:"r" long:"return" description:"Output matched lines"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
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
				))),
		),
	)
}

func createService(opts *logOpts) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.New(sess, aws.NewConfig()), nil
}

type logState struct {
	NextToken *string
	StartTime *int64
}

func (p *awsCloudwatchLogsPlugin) collect() ([]string, error) {
	var nextToken *string
	var startTime *int64
	if s, err := p.loadState(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		if s.StartTime != nil && *s.StartTime > time.Now().Add(-time.Hour).Unix()*1000 {
			nextToken = s.NextToken
			startTime = s.StartTime
		}
	}
	if startTime == nil {
		startTime = aws.Int64(time.Now().Add(-1*time.Minute).Unix() * 1000)
	}
	var messages []string
	for {
		input := &cloudwatchlogs.FilterLogEventsInput{
			StartTime:     startTime,
			LogGroupName:  aws.String(p.LogGroupName),
			NextToken:     nextToken,
			FilterPattern: aws.String(p.Pattern),
		}
		if p.LogStreamNamePrefix != "" {
			input.LogStreamNamePrefix = aws.String(p.LogStreamNamePrefix)
		}
		output, err := p.Service.FilterLogEvents(input)
		if err != nil {
			return nil, err
		}
		for _, event := range output.Events {
			messages = append(messages, *event.Message)
			if startTime == nil || *startTime <= *event.Timestamp {
				startTime = aws.Int64(*event.Timestamp + 1)
			}
		}
		if output.NextToken != nil {
			nextToken = output.NextToken
		}
		if nextToken != nil {
			if err := p.saveState(&logState{nextToken, startTime}); err != nil {
				return nil, err
			}
		}
		if output.NextToken == nil {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	return messages, nil
}

func (p *awsCloudwatchLogsPlugin) loadState() (*logState, error) {
	f, err := os.Open(p.StateFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s logState
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (p *awsCloudwatchLogsPlugin) saveState(s *logState) error {
	err := os.MkdirAll(filepath.Dir(p.StateFile), 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(p.StateFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
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

func (p *awsCloudwatchLogsPlugin) run() *checkers.Checker {
	messages, err := p.collect()
	if err != nil {
		return checkers.Unknown(fmt.Sprint(err))
	}
	return p.check(messages)
}

func run(args []string) *checkers.Checker {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		os.Exit(1)
	}
	p, err := newCloudwatchLogsPlugin(opts, args)
	if err != nil {
		return checkers.Unknown(fmt.Sprint(err))
	}
	return p.run()
}

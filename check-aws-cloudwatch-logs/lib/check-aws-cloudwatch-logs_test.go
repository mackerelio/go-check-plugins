//go:build !windows

package checkawscloudwatchlogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type mockAWSCloudWatchLogsClient struct {
	cloudwatchlogs.FilterLogEventsAPIClient
	outputs []*cloudwatchlogs.FilterLogEventsOutput
}

func (c *mockAWSCloudWatchLogsClient) FilterLogEvents(ctx context.Context, input *cloudwatchlogs.FilterLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var pageNo = 0
	if input.NextToken != nil {
		pageNo, _ = strconv.Atoi(*input.NextToken)
	}
	return c.outputs[pageNo], nil
}

func createMockService() cloudwatchlogs.FilterLogEventsAPIClient {
	return &mockAWSCloudWatchLogsClient{
		outputs: []*cloudwatchlogs.FilterLogEventsOutput{
			{
				NextToken: aws.String("1"),
				Events: []types.FilteredLogEvent{
					{
						EventId:   aws.String("event-id-0"),
						Message:   aws.String("message-0"),
						Timestamp: aws.Int64(0),
					},
					{
						EventId:   aws.String("event-id-1"),
						Message:   aws.String("message-1"),
						Timestamp: aws.Int64(1),
					},
				},
			},
			{
				NextToken: aws.String("2"),
				Events: []types.FilteredLogEvent{
					{
						EventId:   aws.String("event-id-2"),
						Message:   aws.String("message-2"),
						Timestamp: aws.Int64(2),
					},
					{
						EventId:   aws.String("event-id-3"),
						Message:   aws.String("message-3"),
						Timestamp: aws.Int64(3),
					},
					{
						EventId:   aws.String("event-id-4"),
						Message:   aws.String("message-4"),
						Timestamp: aws.Int64(4),
					},
				},
			},
			{
				Events: []types.FilteredLogEvent{
					{
						EventId:   aws.String("event-id-5"),
						Message:   aws.String("message-5"),
						Timestamp: aws.Int64(5),
					},
				},
			},
		},
	}
}

func Test_cloudwatchLogsPlugin_collect(t *testing.T) {
	stateFile := filepath.Join(t.TempDir(), "check-cloudwatch-logs-test-collect")

	p := &awsCloudwatchLogsPlugin{
		Service:   createMockService(),
		StateFile: stateFile,
		logOpts: &logOpts{
			LogGroupName: "test-group",
		},
	}

	t.Run("collect log event messages", func(t *testing.T) {
		messages, err := p.collect(context.Background(), time.Unix(0, 0))
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, len(messages), 6)
		cnt, _ := os.ReadFile(stateFile)
		var s logState
		json.NewDecoder(bytes.NewReader(cnt)).Decode(&s)
		assert.Equal(t, s, logState{StartTime: aws.Int64(5 + 1)})
	})

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		messages, err := p.collect(ctx, time.Unix(0, 0))
		assert.NotEqual(t, err, nil, "err should be someting")
		assert.Equal(t, len(messages), 0)
	})
}

func Test_cloudwatchLogsPlugin_check(t *testing.T) {
	testCases := []struct {
		CriticalOver, WarningOver int
		Pattern                   string
		ReturnContent             bool
		Messages                  []string
		Status                    checkers.Status
		Message                   string
	}{
		{
			CriticalOver: 5,
			WarningOver:  3,
			Pattern:      "Error",
			Messages:     []string{},
			Status:       checkers.OK,
			Message:      "0 messages for pattern /Error/",
		},
		{
			CriticalOver: 5,
			WarningOver:  3,
			Pattern:      "a",
			Messages:     []string{"a0", "a1", "a2"},
			Status:       checkers.OK,
			Message:      "3 messages for pattern /a/",
		},
		{
			CriticalOver: 5,
			WarningOver:  3,
			Pattern:      "a",
			Messages:     []string{"a0", "a1", "a2", "a3", "a4"},
			Status:       checkers.WARNING,
			Message:      "5 > 3 messages for pattern /a/",
		},
		{
			CriticalOver: 5,
			WarningOver:  3,
			Pattern:      "a",
			Messages:     []string{"a0", "a1", "a2", "a3", "a4", "a5"},
			Status:       checkers.CRITICAL,
			Message:      "6 > 5 messages for pattern /a/",
		},
		{
			CriticalOver:  5,
			WarningOver:   3,
			Pattern:       "a",
			ReturnContent: true,
			Messages:      []string{"a0\n", "a1\n", "a2\n", "a3\n", "a4\n", "a5\n"},
			Status:        checkers.CRITICAL,
			Message:       "6 > 5 messages for pattern /a/\na0\na1\na2\na3\na4\na5\n",
		},
	}
	for _, testCase := range testCases {
		p := &awsCloudwatchLogsPlugin{
			logOpts: &logOpts{
				CriticalOver:  testCase.CriticalOver,
				WarningOver:   testCase.WarningOver,
				Pattern:       testCase.Pattern,
				ReturnContent: testCase.ReturnContent,
			},
		}
		res := p.check(testCase.Messages)
		assert.Equal(t, res.Status, testCase.Status)
		assert.Equal(t, res.Message, testCase.Message)
	}
}

func Test_cloudwatchLogsPlugin_options(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want logOpts
	}{
		{
			name: "quoted string",
			args: []string{
				"--log-group-name", `"name"`,
				"--log-stream-name-prefix", `"prefix"`,
				"-p", `"err:"`,
				"-s", `"dir"`,
			},
			want: logOpts{
				LogGroupName:        `"name"`,
				LogStreamNamePrefix: `"prefix"`,
				Pattern:             `"err:"`,
				StateDir:            `"dir"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts logOpts
			_, err := flags.ParseArgs(&opts, tt.args)
			if err != nil {
				t.Fatal("newCloudwatchLogsPlugin:", err)
			}
			assert.Equal(t, tt.want, opts)
		})
	}
}

func Test_createCloudwatchlogsOptions(t *testing.T) {
	tests := []struct {
		opts   *logOpts
		want   cloudwatchlogs.Options
		length int
	}{
		{
			opts:   &logOpts{MaxRetries: 0},
			length: 0,
		},
		{
			opts:   &logOpts{MaxRetries: 1},
			want:   cloudwatchlogs.Options{RetryMaxAttempts: 1},
			length: 1,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			res := createCloudwatchlogsOptions(tt.opts)
			assert.Equal(t, len(res), tt.length)
			if tt.length > 0 {
				opts := cloudwatchlogs.Options{}
				res[0](&opts)
				assert.Equal(t, tt.want, opts)
			}
		})
	}
}

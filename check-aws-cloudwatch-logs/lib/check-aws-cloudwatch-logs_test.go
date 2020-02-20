package checkawscloudwatchlogs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

type mockAWSCloudWatchLogsClient struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	outputs map[string]*cloudwatchlogs.FilterLogEventsOutput
}

func (c *mockAWSCloudWatchLogsClient) FilterLogEvents(input *cloudwatchlogs.FilterLogEventsInput) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	if input.NextToken == nil {
		return c.outputs[""], nil
	}
	if out, ok := c.outputs[*input.NextToken]; ok {
		return out, nil
	}
	return nil, errors.New("invalid NextToken")
}

func createMockService() cloudwatchlogsiface.CloudWatchLogsAPI {
	return &mockAWSCloudWatchLogsClient{
		outputs: map[string]*cloudwatchlogs.FilterLogEventsOutput{
			"": &cloudwatchlogs.FilterLogEventsOutput{
				NextToken: aws.String("1"),
				Events: []*cloudwatchlogs.FilteredLogEvent{
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
			"1": &cloudwatchlogs.FilterLogEventsOutput{
				NextToken: aws.String("2"),
				Events: []*cloudwatchlogs.FilteredLogEvent{
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
			"2": &cloudwatchlogs.FilterLogEventsOutput{
				Events: []*cloudwatchlogs.FilteredLogEvent{
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
	file, _ := ioutil.TempFile("", "check-cloudwatch-logs-test-collect")
	os.Remove(file.Name())
	file.Close()
	defer os.Remove(file.Name())
	p := &awsCloudwatchLogsPlugin{
		Service:   createMockService(),
		StateFile: file.Name(),
		logOpts: &logOpts{
			LogGroupName: "test-group",
		},
	}
	messages, err := p.collect()
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, len(messages), 6)
	cnt, _ := ioutil.ReadFile(file.Name())
	var s logState
	json.NewDecoder(bytes.NewReader(cnt)).Decode(&s)
	assert.Equal(t, *s.NextToken, "2")
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

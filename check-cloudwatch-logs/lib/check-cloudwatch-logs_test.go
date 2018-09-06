package checkcloudwatchlogs

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

type mockCloudWatchLogsClient struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	outputs map[string]*cloudwatchlogs.FilterLogEventsOutput
}

func (c *mockCloudWatchLogsClient) FilterLogEvents(input *cloudwatchlogs.FilterLogEventsInput) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	if input.NextToken == nil {
		return c.outputs[""], nil
	}
	if out, ok := c.outputs[*input.NextToken]; ok {
		return out, nil
	}
	return nil, errors.New("invalid NextToken")
}

func createMockService() cloudwatchlogsiface.CloudWatchLogsAPI {
	return &mockCloudWatchLogsClient{
		outputs: map[string]*cloudwatchlogs.FilterLogEventsOutput{
			"": &cloudwatchlogs.FilterLogEventsOutput{
				NextToken: aws.String("1"),
				Events: []*cloudwatchlogs.FilteredLogEvent{
					{
						EventId: aws.String("event-id-0"),
						Message: aws.String("message-0"),
					},
					{
						EventId: aws.String("event-id-1"),
						Message: aws.String("message-1"),
					},
				},
			},
			"1": &cloudwatchlogs.FilterLogEventsOutput{
				NextToken: aws.String("2"),
				Events: []*cloudwatchlogs.FilteredLogEvent{
					{
						EventId: aws.String("event-id-2"),
						Message: aws.String("message-2"),
					},
					{
						EventId: aws.String("event-id-3"),
						Message: aws.String("message-3"),
					},
					{
						EventId: aws.String("event-id-4"),
						Message: aws.String("message-4"),
					},
				},
			},
			"2": &cloudwatchlogs.FilterLogEventsOutput{
				Events: []*cloudwatchlogs.FilteredLogEvent{
					{
						EventId: aws.String("event-id-5"),
						Message: aws.String("message-5"),
					},
				},
			},
		},
	}
}

func Test_cloudwatchLogsPlugin_run(t *testing.T) {
	p := &cloudwatchLogsPlugin{
		Service:      createMockService(),
		LogGroupName: "test-group",
	}
	messages, err := p.run()
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, len(messages), 6)
}

package processor

import (
	"context"
	"testing"
	"time"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockdb := NewMockDB(ctrl)
	processor := Processor{
		db: mockdb,
	}

	message := []byte(`
	[{
		"absolute_path": "/tmp/file.txt",
		"contents": "hello world",
		"timestamp": "2021-01-01T00:00:00Z"
	}]
	`)

	ctx := context.Background()

	mockdb.EXPECT().InsertFiles(ctx, []lib.File{{
		Absolute_path: "/tmp/file.txt",
		Contents:      "hello world",
		Timestamp:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}}, "USER_ID").Return(nil)

	err := processor.handleMessage(ctx, kafka.Message{Value: message}, "USER_ID")

	require.NoError(t, err)
}

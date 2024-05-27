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
	{
		"user_id": "1",
		"absolute_path": "/tmp/file.txt",
		"contents": "hello world",
		"timestamp": "2021-01-01T00:00:00Z"
	}
	`)

	ctx := context.Background()

	mockdb.EXPECT().InsertFile(ctx, lib.File{
		User_id:       "1",
		Absolute_path: "/tmp/file.txt",
		Contents:      "hello world",
		Timestamp:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}).Return(nil)

	err := processor.handleMessage(ctx, kafka.Message{Value: message})

	require.NoError(t, err)
}

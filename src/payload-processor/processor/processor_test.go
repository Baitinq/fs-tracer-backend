package processor

import (
	"context"
	"testing"

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

	message := []byte("test")

	ctx := context.Background()

	mockdb.EXPECT().TestInsert(ctx, string(message)).Return(nil)

	err := processor.handleMessage(ctx, kafka.Message{Value: message})

	require.NoError(t, err)
}

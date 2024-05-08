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

	mockdb.EXPECT().TestInsert(gomock.Any(), string(message)).Return(nil)

	err := processor.handleMessage(context.Background(), kafka.Message{Value: message})

	require.NoError(t, err)
}

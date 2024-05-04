package processor

import (
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestProcessMessage(t *testing.T) {
	processor := Processor{}

	message := []byte("test")

	err := processor.handleMessage(kafka.Message{Value: message})

	require.NoError(t, err)
}

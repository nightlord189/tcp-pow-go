package protocol

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseMessage(t *testing.T) {
	t.Run("Empty message", func(t *testing.T) {
		input := ""
		msg, err := ParseMessage(input)
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "cannot parse header", err.Error())
	})
	t.Run("Invalid message #2", func(t *testing.T) {
		input := "||"
		msg, err := ParseMessage(input)
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "message doesn't match protocol", err.Error())
	})
	t.Run("Message with non-digit header", func(t *testing.T) {
		input := "test|payload"
		msg, err := ParseMessage(input)
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "cannot parse header", err.Error())
	})
	t.Run("Message with only header", func(t *testing.T) {
		input := "1|"
		msg, err := ParseMessage(input)
		require.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, 1, msg.Header)
		assert.Empty(t, msg.Payload)
	})
	t.Run("Message with header and payload", func(t *testing.T) {
		input := "1|payload"
		msg, err := ParseMessage(input)
		require.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, 1, msg.Header)
		assert.Equal(t, "payload", msg.Payload)
	})
}

package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/config"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/pow"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestProcessRequest(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", &config.Config{HashcashZerosCount: 3})

	t.Run("Quit request", func(t *testing.T) {
		input := fmt.Sprintf("%d|", protocol.Quit)
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, ErrQuit, err)
	})

	t.Run("Invalid request", func(t *testing.T) {
		input := "||"
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "message doesn't match protocol", err.Error())
	})

	t.Run("Unknown header", func(t *testing.T) {
		input := "111|"
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "unknown header", err.Error())
	})

	t.Run("Request challenge", func(t *testing.T) {
		input := fmt.Sprintf("%d|", protocol.RequestChallenge)
		msg, err := ProcessRequest(ctx, input, "client1")
		require.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, protocol.ResponseChallenge, msg.Header)
		//unmarshal msg to check fields
		var hashcash pow.HashcashData
		err = json.Unmarshal([]byte(msg.Payload), &hashcash)
		require.NoError(t, err)
		assert.Equal(t, 3, hashcash.ZerosCount)
		assert.Equal(t, "client1", hashcash.Resource)
		assert.NotEmpty(t, hashcash.Rand)
	})

	t.Run("Request resource without solution", func(t *testing.T) {
		input := fmt.Sprintf("%d|", protocol.RequestResource)
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.True(t, strings.Contains(err.Error(), "err unmarshal hashcash"))
	})

	t.Run("Request resource with wrong resource", func(t *testing.T) {
		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: 4,
			Date:       time.Now().Unix(),
			Resource:   "client2",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    100,
		}
		marshaled, err := json.Marshal(hashcash)
		require.NoError(t, err)
		input := fmt.Sprintf("%d|%s", protocol.RequestResource, string(marshaled))
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "invalid hashcash resource", err.Error())
	})

	t.Run("Request resource with invalid solution and 0 counter", func(t *testing.T) {
		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: 10,
			Date:       time.Now().Unix(),
			Resource:   "client1",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    0,
		}
		marshaled, err := json.Marshal(hashcash)
		require.NoError(t, err)
		input := fmt.Sprintf("%d|%s", protocol.RequestResource, string(marshaled))
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "invalid hashcash", err.Error())
	})

	t.Run("Request resource with invalid solution", func(t *testing.T) {
		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: 10,
			Date:       time.Now().Unix(),
			Resource:   "client1",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    100,
		}
		marshaled, err := json.Marshal(hashcash)
		require.NoError(t, err)
		input := fmt.Sprintf("%d|%s", protocol.RequestResource, string(marshaled))
		msg, err := ProcessRequest(ctx, input, "client1")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Equal(t, "invalid hashcash", err.Error())
	})

	t.Run("Request resource with correct solution", func(t *testing.T) {
		date := time.Date(2022, 3, 13, 2, 30, 0, 0, time.UTC)
		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: 3,
			Date:       date.Unix(),
			Resource:   "client1",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    5001,
		}
		marshaled, err := json.Marshal(hashcash)
		require.NoError(t, err)
		input := fmt.Sprintf("%d|%s", protocol.RequestResource, string(marshaled))
		msg, err := ProcessRequest(ctx, input, "client1")
		require.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Contains(t, Quotes, msg.Payload)
	})
}
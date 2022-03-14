package client

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

type MockConnection struct {
	ReadFunc  func([]byte) (int, error)
	WriteFunc func([]byte) (int, error)
}

func (m MockConnection) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

func (m MockConnection) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func TestHandleConnection(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", &config.Config{HashcashMaxIterations: 1000000})

	t.Run("Write error", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, fmt.Errorf("test write error")
			},
		}
		_, err := HandleConnection(ctx, mock, mock)
		assert.Error(t, err)
		assert.Equal(t, "err send request: test write error", err.Error())
	})

	t.Run("Read error", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				return 0, fmt.Errorf("test read error")
			},
		}
		_, err := HandleConnection(ctx, mock, mock)
		assert.Error(t, err)
		assert.Equal(t, "err read msg: test read error", err.Error())
	})

	t.Run("Read response in bad format", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				return fillTestReadBytes("||\n", p), nil
			},
		}
		_, err := HandleConnection(ctx, mock, mock)
		assert.Error(t, err)
		assert.Equal(t, "err parse msg: message doesn't match protocol", err.Error())
	})

	t.Run("Read response with hashcash in bad format", func(t *testing.T) {
		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				return fillTestReadBytes(fmt.Sprintf("%d|{wrong_json}\n", protocol.ResponseChallenge), p), nil
			},
		}
		_, err := HandleConnection(ctx, mock, mock)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "err parse hashcash"))
	})

	t.Run("Success", func(t *testing.T) {
		date := time.Date(2022, 3, 13, 2, 30, 0, 0, time.UTC)
		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: 3,
			Date:       date.Unix(),
			Resource:   "client1",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    0,
		}

		// counter for reading attempts to change content
		readAttempt := 0

		writeAttempt := 0

		mock := MockConnection{
			WriteFunc: func(p []byte) (int, error) {
				if writeAttempt == 0 {
					writeAttempt++
					assert.Equal(t, "1|\n", string(p))
				} else {
					msg, err := protocol.ParseMessage(string(p))
					require.NoError(t, err)
					var writtenHashcash pow.HashcashData
					err = json.Unmarshal([]byte(msg.Payload), &writtenHashcash)
					require.NoError(t, err)
					// checking that counter increased
					assert.Equal(t, 5001, writtenHashcash.Counter)
					_, err = writtenHashcash.ComputeHashcash(0)
					assert.NoError(t, err)
				}
				return 0, nil
			},
			ReadFunc: func(p []byte) (int, error) {
				if readAttempt == 0 {
					marshaled, err := json.Marshal(hashcash)
					require.NoError(t, err)
					readAttempt++
					return fillTestReadBytes(fmt.Sprintf("%d|%s\n", protocol.ResponseChallenge, string(marshaled)), p), nil
				} else {
					// second read, send quote
					return fillTestReadBytes(fmt.Sprintf("%d|test quote\n", protocol.ResponseChallenge), p), nil
				}
			},
		}
		response, err := HandleConnection(ctx, mock, mock)
		assert.NoError(t, err)
		assert.Equal(t, "test quote", response)
	})
}

// fillTestReadBytes - helper to easier mock Reader
func fillTestReadBytes(str string, p []byte) int {
	dataBytes := []byte(str)
	counter := 0
	for i := range dataBytes {
		p[i] = dataBytes[i]
		counter++
		if counter >= len(p) {
			break
		}
	}
	return counter
}

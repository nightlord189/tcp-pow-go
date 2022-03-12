package pkg

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestCase struct {
	input    interface{}
	expected interface{}
}

func TestHash(t *testing.T) {
	result := sha1Hash("testdatalong 1231378612")
	assert.Equal(t, "26e48dc4c6fd473c87e9c4928d8f4bd45f508603", result)
	result = sha1Hash("super800")
	assert.Equal(t, "a6a735584d2a32fbbd5af4cb6d9931167d7bb2db", result)
}

func TestIsHashCorrect(t *testing.T) {
	result := IsHashCorrect("0000cd3d39d4c11079d167e870dbc5873f3c9169", 4)
	assert.True(t, result)
	result = IsHashCorrect("0000", 5)
	assert.False(t, result)
	result = IsHashCorrect("00018f38327b85110de794711aa02926ae8f7f76", 4)
	assert.False(t, result)
}

func TestComputeHashcash(t *testing.T) {
	t.Parallel()

	t.Run("4 zeros", func(t *testing.T) {
		date := time.Date(2022, 3, 13, 2, 28, 0, 0, time.UTC)
		input := HashcashData{
			Version:    1,
			ZerosCount: 4,
			Date:       date.Unix(),
			Resource:   "some_useful_data",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123459))),
			Counter:    0,
		}
		result, err := input.ComputeHashcash(-1)
		require.NoError(t, err)
		assert.Equal(t, 26394, result.Counter)
	})
	t.Run("5 zeros", func(t *testing.T) {
		date := time.Date(2022, 3, 13, 2, 30, 0, 0, time.UTC)
		input := HashcashData{
			Version:    1,
			ZerosCount: 5,
			Date:       date.Unix(),
			Resource:   "some_useful_data",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    0,
		}
		result, err := input.ComputeHashcash(-1)
		require.NoError(t, err)
		assert.Equal(t, 36258, result.Counter)
	})
	t.Run("impossible challenge", func(t *testing.T) {
		date := time.Date(2022, 3, 13, 2, 30, 0, 0, time.UTC)
		input := HashcashData{
			Version:    1,
			ZerosCount: 10,
			Date:       date.Unix(),
			Resource:   "some_useful_data",
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
			Counter:    0,
		}
		result, err := input.ComputeHashcash(10)
		require.Error(t, err)
		assert.Equal(t, 11, result.Counter)
		assert.Equal(t, "max iterations exceeded", err.Error())
	})
}

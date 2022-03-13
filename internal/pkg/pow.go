package pkg

import (
	"crypto/sha1"
	"fmt"
)

const zeroByte = 48

type HashcashData struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
}

func (h HashcashData) Stringify() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

func (h HashcashData) ComputeHashcash(maxIterations int) (HashcashData, error) {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		header := h.Stringify()
		hash := sha1Hash(header)
		//fmt.Println(header, hash)
		if IsHashCorrect(hash, h.ZerosCount) {
			return h, nil
		}
		h.Counter++
	}
	return h, fmt.Errorf("max iterations exceeded")
}

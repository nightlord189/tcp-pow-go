package pow

import (
	"crypto/sha1"
	"fmt"
)

// to avoid magic comparisons
const zeroByte = 48

// HashcashData - struct with fields of Hashcash
// https://en.wikipedia.org/wiki/Hashcash
type HashcashData struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
}

// Stringify - stringifies hashcash for next sending it on TCP
func (h HashcashData) Stringify() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

// sha1Hash - calculates sha1 hash from given string
func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// IsHashCorrect - checks that hash has leading <zerosCount> zeros
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

// ComputeHashcash - calculates correct hashcash by bruteforce
// until the resulting hash satisfies the condition of IsHashCorrect
// maxIterations to prevent endless computing (0 or -1 to disable it)
func (h HashcashData) ComputeHashcash(maxIterations int) (HashcashData, error) {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		header := h.Stringify()
		hash := sha1Hash(header)
		//fmt.Println(header, hash)
		if IsHashCorrect(hash, h.ZerosCount) {
			return h, nil
		}
		// if hash don't have needed count of leading zeros, we are increasing counter and try next hash
		h.Counter++
	}
	return h, fmt.Errorf("max iterations exceeded")
}

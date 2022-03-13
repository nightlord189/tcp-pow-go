package main

import (
	"encoding/base64"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg"
	"time"
)

func main() {
	fmt.Println("start")

	date := time.Date(2022, 3, 13, 2, 30, 0, 0, time.UTC)
	hashcash := pkg.HashcashData{
		Version:    1,
		ZerosCount: 3,
		Date:       date.Unix(),
		Resource:   "client1",
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
		Counter:    100,
	}
	hashcash, _ = hashcash.ComputeHashcash(-1)
	fmt.Println(hashcash)
}

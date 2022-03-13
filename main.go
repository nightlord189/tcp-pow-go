package main

import (
	"encoding/base64"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg"
	"time"
)

func main() {
	fmt.Println("start")

	date := time.Now()

	date = time.Date(2022, 3, 13, 2, 30, 0, 0, time.UTC)

	dd := pkg.HashcashData{
		Version:    1,
		ZerosCount: 5,
		Date:       date.Unix(),
		Resource:   "some_useful_data",
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
		Counter:    0,
	}
	dd, _ = dd.ComputeHashcash(-1)
	fmt.Println(dd)

	msg := pkg.Message{
		Header: pkg.RequestChallenge,
		Data:   "ddd",
	}
	fmt.Println(msg.Stringify())
}

package server

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var Quotes = []string{
	"All saints who remember to keep and do these sayings, " +
		"walking in obedience to the commandments, " +
		"shall receive health in their navel and marrow to their bones",

	"And shall find wisdom and great treasures of knowledge, even hidden treasures",

	"And shall run and not be weary, and shall walk and not faint",

	"And I, the Lord, give unto them a promise, " +
		"that the destroying angel shall pass by them, " +
		"as the children of Israel, and not slay them",
}

func HandleConnection(conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')

		switch err {
		case nil:
			handleRequest(req, conn)
			req := strings.TrimSpace(req)
			if req == ":QUIT" {
				fmt.Println("client requested server to close the connection so closing")
				return
			} else {
				fmt.Println("msg received:", req)
			}
		case io.EOF:
			fmt.Println("client closed the connection by terminating the process")
			return
		default:
			fmt.Printf("error: %v\n", err)
			return
		}

		_, err = conn.Write([]byte("response\n"))
		if err != nil {
			fmt.Printf("failed to respond to client: %v\n", err)
		} else {
			fmt.Println("msg sent: response")
		}
	}
}

func handleRequest(msg string, conn net.Conn) bool {
	msg = strings.TrimSpace(msg)
	var msgType int
	parts := strings.Split(msg, "|")
	if len(parts) != 0 {
		fmt.Printf("msg %s doesn't match protocol\n", msg)
		return true
	}
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("cannot parse header %s\n", parts[0])
		return true
	}
	switch msgType {
	case pkg.Quit:
		fmt.Printf("client %s requests quit\n", conn.RemoteAddr())
		return false
	case pkg.RequestChallenge:
		fmt.Printf("client %s requests challenge\n", conn.RemoteAddr())
		msg := pkg.Message{
			Header: pkg.ResponseChallenge,
		}
		date := time.Now()
		hashcash := pkg.HashcashData{
			Version:    1,
			ZerosCount: 5,
			Date:       date.Unix(),
			Resource:   conn.RemoteAddr().String(),
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(100000)))),
			Counter:    0,
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			fmt.Printf("err marshal hashcash: %v\n", err)
			msg.Data = "error marshal hashcash to json"
			sendMsg(msg, conn)
			return true
		}
		msg.Data = string(hashcashMarshaled)
		sendMsg(msg, conn)
		return true
	case pkg.RequestResource:
		msg := pkg.Message{
			Header: pkg.ResponseResource,
		}
		//parse client's solution
		var hashcash pkg.HashcashData
		err := json.Unmarshal([]byte(parts[1]), &hashcash)
		if err != nil {
			msg.Data = "error parse hashcash json"
			sendMsg(msg, conn)
			return true
		}
		_, err = hashcash.ComputeHashcash(hashcash.Counter)
		if err != nil {
			msg.Data = "error verify hashcash"
			sendMsg(msg, conn)
			return true
		}
		//get random quote
		msg.Data = Quotes[rand.Intn(4)]
		return true
	default:
		fmt.Printf("unknown header %d\n", msgType)
		return false
	}
}

func sendMsg(msg pkg.Message, conn net.Conn) {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	if err != nil {
		fmt.Printf("err send message to %s: %v\n", conn.RemoteAddr(), err)
	}
}

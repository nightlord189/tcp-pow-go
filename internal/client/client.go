package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg"
	"net"
)

func Run(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	fmt.Println("connected")
	defer conn.Close()

	err = HandleConnection(conn)
	if err != nil {
		return err
	}
	fmt.Println("stop client")
	return nil
}

// HandleConnection - scenario for TCP-client
// 1. request challenge from server
// 2. compute hashcash to check Proof of Work
// 3. send hashcash solution back to server
// 4. get result quote from server
func HandleConnection(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	// 1. requesting challenge
	err := sendMsg(pkg.Message{
		Header: pkg.RequestChallenge,
	}, conn)
	if err != nil {
		return fmt.Errorf("err send request: %w", err)
	}

	// reading and parsing response
	msgStr, err := readConnMsg(reader)
	if err != nil {
		return fmt.Errorf("err read msg: %w", err)
	}
	msg, err := pkg.ParseMessage(msgStr)
	if err != nil {
		return fmt.Errorf("err parse msg: %w", err)
	}
	var hashcash pkg.HashcashData
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return fmt.Errorf("err parse hashcash: %w", err)
	}
	fmt.Println("got hashcash:", hashcash)

	// 2. got challenge, compute hashcash
	hashcash, err = hashcash.ComputeHashcash(1000000)
	if err != nil {
		return fmt.Errorf("err compute hashcash: %w", err)
	}
	fmt.Println("hashcash computed:", hashcash)
	// marshal solution to json
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return fmt.Errorf("err marshal hashcash: %w", err)
	}

	// 3. send challenge solution back to server
	err = sendMsg(pkg.Message{
		Header:  pkg.RequestResource,
		Payload: string(byteData),
	}, conn)
	if err != nil {
		return fmt.Errorf("err send request: %w", err)
	}
	fmt.Println("challenge sent to server")

	// 4. get result quote from server
	msgStr, err = readConnMsg(reader)
	if err != nil {
		return fmt.Errorf("err read msg: %w", err)
	}
	msg, err = pkg.ParseMessage(msgStr)
	if err != nil {
		return fmt.Errorf("err parse msg: %w", err)
	}

	fmt.Println("quote result", msg.Payload)

	return nil
}

func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

func sendMsg(msg pkg.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}

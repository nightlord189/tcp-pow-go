package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

// Header of TCP-message in protocol, means type of message
const (
	Quit              = iota // on quit each side (server or client) should close connection
	RequestChallenge         // from client to server - request new challenge from server
	ResponseChallenge        // from server to client - message with challenge for client
	RequestResource          // from client to server - message with solved challenge
	ResponseResource         // from server to client - message with useful info is solution is correct, or with error if not
)

// Message - message struct for both server and client
type Message struct {
	Header  int    //type of message
	Payload string //payload, could be json, quote or be empty
}

// Stringify - stringify message to send it by tcp-connection
// divider between header and payload is |
func (m *Message) Stringify() string {
	return fmt.Sprintf("%d|%s", m.Header, m.Payload)
}

// ParseMessage - parses Message from str, checks header and payload
func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	var msgType int
	// message has view as 1|payload (payload is optional)
	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 { //only 1 or 2 parts allowed
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	// try to parse header
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse header")
	}
	msg := Message{
		Header: msgType,
	}
	// last part after | is payload
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}

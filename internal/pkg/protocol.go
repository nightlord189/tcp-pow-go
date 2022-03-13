package pkg

import "fmt"

//Header of TCP-message in protocol, means type of message
const (
	Quit              = iota //on quit each side (server or client) should close connection
	RequestChallenge         //from client to server - request new challenge from server
	ResponseChallenge        //from server to client - message with challenge for client
	RequestResource          //from client to server - message with solved challenge
	ResponseResource         //from server to client - message with useful info is solution is correct, or with error if not
)

type Message struct {
	Header int
	Data   string
}

func (m *Message) Stringify() string {
	return fmt.Sprintf("%d|%s", m.Header, m.Data)
}

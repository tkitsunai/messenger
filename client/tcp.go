package client

import (
	"encoding/json"
	"fmt"
	"io"
	"message/core"
	"net"
)

type Tcp struct {
	conn     io.ReadWriteCloser
	encoder  *json.Encoder
	host     string
	messages chan core.Message
}

func NewTcp(host string) (*Tcp, error) {
	client := Tcp{
		host:     host,
		messages: make(chan core.Message),
	}
	return &client, client.connect()
}

func (t *Tcp) connect() error {
	conn, err := net.Dial("tcp", t.host)
	if err != nil {
		return err
	}
	t.conn = conn
	t.encoder = json.NewEncoder(conn)

	decoder := json.NewDecoder(conn)
	go func() {
		for {
			msg := core.Message{}
			if err := decoder.Decode(&msg); err != nil {
				conn.Close()
				close(t.messages)
				return
			}
			t.messages <- msg
		}
	}()
	return nil
}

func (t *Tcp) Close() error {
	return t.conn.Close()
}

func (t *Tcp) Subscribe(subjects []string) error {
	if len(subjects) == 0 {
		return fmt.Errorf("unable to subscribe - missing subject")
	}
	return t.encoder.Encode(&core.Message{
		Command:  core.CmdSubscribe,
		Subjects: subjects,
	})
}

func (t *Tcp) Publish(subjects []string, payload string) error {
	if len(subjects) == 0 {
		return fmt.Errorf("unable to subscribe - missing subject")
	}
	if len(payload) == 0 {
		return fmt.Errorf("unable to publish - missing payload")
	}

	return t.encoder.Encode(&core.Message{
		Command:  core.CmdPublish,
		Subjects: subjects,
		Payload:  payload,
	})
}

func (t *Tcp) Messages() <-chan core.Message {
	return t.messages
}

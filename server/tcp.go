package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tkitsunai/messenger/core"
	"github.com/tkitsunai/messenger/utils"
	"io"
	"net"
)

type TcpServer struct {
	StartHandler core.ServerStartHandler
}

func NewTcpServer() *TcpServer {
	return &TcpServer{
		StartHandler: StartTcp,
	}
}

func StartTcp(uri string, errCh chan<- error) {
	utils.GetLogger().Info().Msg("welcome to messenger")

	ln, err := net.Listen("tcp", uri)

	if err != nil {
		errCh <- err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errCh <- fmt.Errorf("failed to accept tcp connection: %s", err.Error())
			}
			go handleConnection(conn, errCh)
		}
	}()
}

func handleConnection(conn net.Conn, errCh chan<- error) {
	broker := core.NewBroker()
	defer broker.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	go func() {
		for msg := range broker.Pipe {
			if err := encoder.Encode(msg); err != nil {
				errCh <- fmt.Errorf("Failed to pubilsh broker.Pipe contents to clients - %s ", err.Error())
				break
			}
		}
	}()

	handlers := getAllCommandHandlers()

	for {
		msg := core.Message{}
		if err := decoder.Decode(&msg); err != nil {
			switch {
			case err == io.EOF:
				// disconnected?
			case errors.Is(err, io.ErrUnexpectedEOF):
			default:
				errCh <- fmt.Errorf("failed to decode message from client %s", err.Error())
			}
			return
		}

		handler, found := handlers[msg.Command]
		if !found {
			err := encoder.Encode(core.Message{
				Error: "Unknown Command",
			})
			if err != nil {
				utils.GetLogger().Error().Msg(fmt.Sprintf("failed encoded: %s", err))
			}
			continue
		}

		if err := handler(broker, msg); err != nil {
			err := encoder.Encode(core.Message{
				Command: msg.Command,
				Error:   err.Error(),
			})
			if err != nil {
				utils.GetLogger().Error().Msg(fmt.Sprintf("failed encoded: %s", err))
			}
			continue
		}
	}
}

func getAllCommandHandlers() map[string]core.HandlerFn {
	return map[string]core.HandlerFn{
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnSubscribe,
		"publish":     handlePublish,
	}
}

func handleSubscribe(broker *core.Broker, msg core.Message) error {
	return broker.Subscribe(msg.Subjects)
}

func handleUnSubscribe(broker *core.Broker, msg core.Message) error {
	return broker.UnSubscribe(msg.Subjects)
}

func handlePublish(broker *core.Broker, msg core.Message) error {
	return broker.Publish(msg.Subjects, msg.Payload)
}

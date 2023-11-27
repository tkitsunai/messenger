package core

import "github.com/google/uuid"

// Client
// TODO Incomplete implementation.
type Client struct {
	// ユーザID、セッションIDを持つ
	userId    string
	sessionId string

	// Nodeへの参照を持つ
	node *Node
}

type ClientCloseFunc func() error

// NewClient
// TODO Incomplete implementation.
func NewClient(broker *MemoryBroker, node *Node) (*Client, ClientCloseFunc, error) {
	userId, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}

	sessionId, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}

	c := &Client{
		userId:    userId.String(),
		sessionId: sessionId.String(),
		node:      node,
	}

	return c, func() error {
		return c.Disconnect()
	}, nil
}

func (c *Client) UserID() string {
	return c.userId
}

func (c *Client) Disconnect() error {
	// TODO
	panic("TODO implementation.")
}

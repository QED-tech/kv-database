package network

import (
	"fmt"
	"net"
)

type Client struct {
	address string
}

func NewTCPClient(address string) *Client {
	return &Client{
		address: address,
	}
}

func (c *Client) Send(query string) (string, error) {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return "", fmt.Errorf("failed to connect to database: %v", err)
	}

	if _, err = conn.Write([]byte(query)); err != nil {
		return "", fmt.Errorf("failed to write query: %v", err)
	}

	buf := make([]byte, 1024)

	if _, err = conn.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read response from database: %v", err)
	}

	return string(buf), nil
}

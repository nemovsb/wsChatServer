package app

import "github.com/gorilla/websocket"

type Connect struct {
	*Client
	ws websocket.Conn
}

func NewConnect(client *Client, ws websocket.Conn) *Connect {
	return &Connect{
		Client: client,
		ws:     ws,
	}
}

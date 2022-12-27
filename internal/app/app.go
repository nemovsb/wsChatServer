package app

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type App struct {
	Chanels map[string]Chanel
}

func NewApp() *App {

	ch := make(map[string]Chanel)

	return &App{
		Chanels: ch,
	}
}

// CloseConnections
func (a *App) CloseConnections() {
	for _, ch := range a.Chanels {
		for _, conn := range ch.Connects {
			err := conn.ws.Close()
			if err != nil {
				log.Println("close connect error: ", err)
			}
			ch.DelConn(&conn)
		}
	}
}

// Serve served websocket connections
func (a *App) Serve(ws websocket.Conn) {

	handshakeMessage, err := read(&ws)
	if err != nil {
		fmt.Printf("read error: %s\n", err)
		return
	}

	fmt.Printf("handshakeMessage: %+v\n", handshakeMessage)

	client := NewClient(handshakeMessage.ClientName)
	connect := NewConnect(client, ws)

	ch, ok := a.Chanels[handshakeMessage.ChanelName]

	if !ok {
		ch = NewChanel(*connect, handshakeMessage)
		a.Chanels[handshakeMessage.ChanelName] = ch

	} else {
		ch.AddConn(connect)
	}

	for *(ch.ClientsCount) != 0 {

		message, err := read(&ws)
		if err != nil {
			ch.DelConn(connect)
			fmt.Printf("read error: %s\n", err)

			if *(ch.ClientsCount) == 0 {
				fmt.Printf("chanel %s delete\n", message.ChanelName)
				delete(a.Chanels, message.ChanelName)

			}

			return
		}

		_, ok := a.Chanels[message.ChanelName]
		if ok {

			for _, conn := range a.Chanels[message.ChanelName].Connects {
				if conn.Nickname == client.Nickname {
					continue
				}

				err := send(conn, message)
				if err != nil {
					ch.DelConn(&conn)
					fmt.Printf("send error: %s\n", err)
					continue
				}
			}
		}

	}

}

func read(conn *websocket.Conn) (Message, error) {
	for {

		message := Message{}

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
		}

		return message, err

	}
}

func send(c Connect, m Message) error {

	fmt.Printf("send message: %+v\n", m)

	err := c.ws.WriteJSON(&m)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

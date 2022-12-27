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
	//fmt.Printf("ch.ClientsCount	get: %d\n", *(ch.ClientsCount))

	if !ok {

		fmt.Printf("!ok: %v\n", !ok)
		ch = NewChanel(*connect, handshakeMessage)
		a.Chanels[handshakeMessage.ChanelName] = ch
		fmt.Printf("ch.ClientsCount	NewChanel: %d\n", *(ch.ClientsCount))
	} else {
		fmt.Printf("ok: %v\n", ok)
		ch.AddConn(connect)
		fmt.Printf("ch.ClientsCount	AddConn: %d\n", *(ch.ClientsCount))
	}

	//fmt.Printf("ch.Connects: %v\n", ch.Connects)

	for *(ch.ClientsCount) != 0 {

		message, err := read(&ws)
		if err != nil {
			ch.DelConn(connect)
			fmt.Printf("read error: %s\n", err)

			if *(ch.ClientsCount) == 0 {
				fmt.Printf("chanel %s delete\n", ch.Name)
				delete(a.Chanels, ch.Name)

			}

			return
		}

		fmt.Printf("message: %+v\n", message)

		for _, conn := range ch.Connects {
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

		fmt.Printf("ch.ClientsCount: %d\n", *(ch.ClientsCount))
	}

}

func read(conn *websocket.Conn) (Message, error) {
	for {

		message := Message{}

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return Message{}, err
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

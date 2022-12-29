package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type App struct {
	ServeChannels chan *Chanel
	Chanels       map[string]*Chanel
}

func NewApp() *App {
	return &App{
		ServeChannels: make(chan *Chanel),
		Chanels:       make(map[string]*Chanel),
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

func (a *App) ServeChanels() {
	var wg sync.WaitGroup

	for ch := range a.ServeChannels {
		wg.Add(1)
		go func(channel *Chanel) {
			channel.Broadcast()
			delete(a.Chanels, channel.Name)
			wg.Done()
		}(ch)
	}
	wg.Wait()
}

// Serve served websocket connections
func (a *App) Serve(ws websocket.Conn) {

	fmt.Printf("new connect! \n")

	var wg sync.WaitGroup

	inMessage := make(chan Message)

	//read incoming messages
	wg.Add(1)
	readErr := make(chan error)
	go func() {
		readErr <- read(&ws, inMessage)

		wg.Done()

	}()

	handshakeMessage := <-inMessage

	client := NewClient(handshakeMessage.ClientName)
	connect := NewConnect(client, ws)

	_, ok := a.Chanels[handshakeMessage.ChanelName]

	if !ok {
		a.Chanels[handshakeMessage.ChanelName] = NewChanel(*connect, handshakeMessage)
		a.ServeChannels <- a.Chanels[handshakeMessage.ChanelName]
		fmt.Printf("create new channel! owner: %s\n", handshakeMessage.ClientName)

		// wg.Add(1)
		// go func() {
		// 	ch := a.Chanels[handshakeMessage.ChanelName]
		// 	ch.Broadcast()

		// 	wg.Done()
		// }()

	} else {
		a.Chanels[handshakeMessage.ChanelName].AddConn(connect)
	}

	for {
		select {
		case <-readErr:
			{
				if a.Chanels[handshakeMessage.ChanelName].DelConn(connect) {
					a.Chanels[handshakeMessage.ChanelName].cancelCtx()
				}
				return
			}
		case message := <-inMessage:
			{

				if message.ChanelName == handshakeMessage.ChanelName {

					fmt.Printf("send message to broadcast %+v, %v\n", message, &(a.Chanels[handshakeMessage.ChanelName].in))

					a.Chanels[handshakeMessage.ChanelName].in <- message
				} else {
					(fmt.Printf("wrong channel %s\n", message.ChanelName))
				}
			}
		}
	}

	wg.Wait()

}

func read(conn *websocket.Conn, in chan<- Message) error {
	for {

		message := Message{}

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Printf("new message %+v\n", message)

		in <- message

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

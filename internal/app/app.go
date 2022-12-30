package app

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

type Connector interface {
	GetMessages() (in <-chan Message, event <-chan *Event, err error)
	SendMessage(message Message, connect Connect) error

	// AddConnection()
	// DelConnection()
}

type App struct {
	Connector
	SendTasks     chan SendTask
	ServeChannels chan *Chanel
	Chanels       map[string]*Chanel
}

func NewApp(c Connector) *App {
	return &App{
		Connector:     c,
		SendTasks:     make(chan SendTask),
		ServeChannels: make(chan *Chanel),
		Chanels:       make(map[string]*Chanel),
	}
}

func (a *App) Start(ctx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		a.ServeChanels()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a.ServeConnections()
		wg.Done()
	}()

}

func (a *App) ServeChanels() {
	var wg sync.WaitGroup

	for ch := range a.ServeChannels {
		wg.Add(1)

		go func(channel *Chanel) {
			channel.Broadcast(a.SendTasks)
			delete(a.Chanels, channel.Name)
			wg.Done()
		}(ch)
	}

	wg.Wait()
}

func (a *App) DeleteConnect(id int) {
	for _, ch := range a.Chanels {
		for _, conn := range ch.Connects {
			if conn.Id == id {
				ch.DelConn(&conn)
				return
			}
		}
	}
}

func (a *App) ServeConnections() {
	var wg sync.WaitGroup
	defer wg.Wait()

	inMessage, events, err := a.Connector.GetMessages()
	if err != nil {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			event := <-events

			switch event.Type {
			case EmptyEvent:
				{
					continue
				}

			case DeleleConnectionEvent:
				{
					id, err := strconv.Atoi(event.Data)
					if err != nil {
						fmt.Printf("get id error: %s\n", err)
					}

					a.DeleteConnect(id)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			message := <-inMessage

			if chat, ok := a.Chanels[message.ChanelName]; !ok {
				a.CreateChannel(message)

				fmt.Printf("create channel:\n%+v\n", a.Chanels[message.ChanelName])

			} else if connect, ok := chat.Connects[message.ClientName]; !ok {
				connect = *NewConnect(message.FromConnId, NewClient(message.ClientName))
				chat.AddConn(&connect)

			} else {
				a.Chanels[message.ChanelName].in <- message
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			task := <-a.SendTasks
			err := a.Connector.SendMessage(task.M, task.To)
			if err != nil {
				connect := a.Chanels[task.M.ChanelName].Connects[task.M.ClientName]
				if a.Chanels[task.M.ChanelName].DelConn(&connect) {
					a.Chanels[task.M.ChanelName].cancelCtx()
					delete(a.Chanels, task.M.ChanelName)
				}
			}
		}
	}()

}

func (a *App) CreateChannel(m Message) {
	connect := NewConnect(m.FromConnId, NewClient(m.ClientName))
	a.Chanels[m.ChanelName] = NewChanel(*connect, m)
}

// Serve served websocket connections
// func (a *App) Serve(ws websocket.Conn) {

// 	fmt.Printf("new connect! \n")

// 	var wg sync.WaitGroup

// 	inMessage := make(chan Message)

// 	//read incoming messages
// 	wg.Add(1)
// 	readErr := make(chan error)
// 	go func() {
// 		readErr <- read(&ws, inMessage)

// 		wg.Done()

// 	}()

// 	handshakeMessage := <-inMessage

// 	client := NewClient(handshakeMessage.ClientName)
// 	connect := NewConnect(client, ws)

// 	_, ok := a.Chanels[handshakeMessage.ChanelName]

// 	if !ok {
// 		a.Chanels[handshakeMessage.ChanelName] = NewChanel(*connect, handshakeMessage)
// 		a.ServeChannels <- a.Chanels[handshakeMessage.ChanelName]
// 		fmt.Printf("create new channel! owner: %s\n", handshakeMessage.ClientName)

// 		// wg.Add(1)
// 		// go func() {
// 		// 	ch := a.Chanels[handshakeMessage.ChanelName]
// 		// 	ch.Broadcast()

// 		// 	wg.Done()
// 		// }()

// 	} else {
// 		a.Chanels[handshakeMessage.ChanelName].AddConn(connect)
// 	}

// 	for {
// 		select {
// 		case <-readErr:
// 			{
// 				if a.Chanels[handshakeMessage.ChanelName].DelConn(connect) {
// 					a.Chanels[handshakeMessage.ChanelName].cancelCtx()
// 				}
// 				return
// 			}
// 		case message := <-inMessage:
// 			{

// 				if message.ChanelName == handshakeMessage.ChanelName {

// 					fmt.Printf("send message to broadcast %+v, %v\n", message, &(a.Chanels[handshakeMessage.ChanelName].in))

// 					a.Chanels[handshakeMessage.ChanelName].in <- message
// 				} else {
// 					(fmt.Printf("wrong channel %s\n", message.ChanelName))
// 				}
// 			}
// 		}
// 	}

// 	wg.Wait()

// }

func (a *App) CloseConnections() {
	for _, ch := range a.Chanels {
		for _, conn := range ch.Connects {
			ch.DelConn(&conn)
		}
	}
}

// func read(conn *websocket.Conn, in chan<- Message) error {
// 	for {

// 		message := Message{}

// 		err := conn.ReadJSON(&message)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}

// 		fmt.Printf("new message %+v\n", message)

// 		in <- message

// 	}
// }

// func send(c Connect, m Message) error {

// 	fmt.Printf("send message: %+v\n", m)

// 	err := c.ws.WriteJSON(&m)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	return err
// }

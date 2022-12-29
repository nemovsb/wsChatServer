package ws

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"ws_server/internal/app"

	"github.com/gorilla/websocket"
)

var id = 1

func getNewId() int {
	id++
	return id
}

type Connect struct {
	Id int

	Client app.Client
	Chat   app.Chanel
	wsConn *websocket.Conn
}

func NewConnect(ws *websocket.Conn) *Connect {
	return &Connect{
		Id:     getNewId(),
		Client: app.Client{},
		Chat:   app.Chanel{},
		wsConn: ws,
	}
}

type WebSocketConnector struct {
	InputMessage chan app.Message
	Events       chan *app.Event

	InputConn chan *websocket.Conn
	Connects  map[int]*Connect
}

func NewConnector() *WebSocketConnector {
	return &WebSocketConnector{
		InputMessage: make(chan app.Message),
		Events:       make(chan *app.Event),
		InputConn:    make(chan *websocket.Conn),
		Connects:     make(map[int]*Connect),
	}
}

func (ws *WebSocketConnector) AddConn(wsConnect *websocket.Conn) *Connect {
	ws.Connects[id] = NewConnect(wsConnect)
	return ws.Connects[id]
}

func (ws *WebSocketConnector) DelConn(id int) {
	delete(ws.Connects, id)
}

func (ws *WebSocketConnector) Start() {
	var wg sync.WaitGroup
	defer wg.Wait()

	for wsConn := range ws.InputConn {

		wg.Add(1)
		go func(wsc *websocket.Conn) {
			wg.Done()

			connect := ws.AddConn(wsc)
			err := read(connect.wsConn, ws.InputMessage)
			if err != nil {
				fmt.Printf("read error: %s", err)
			}
			ws.DelConn(connect.Id)
			delEvent := app.NewEvent()
			delEvent.Type = app.DeleleConnectionEvent
			delEvent.Data = strconv.Itoa(connect.Id)

			ws.Events <- delEvent

		}(wsConn)
	}
}

//func

func read(conn *websocket.Conn, in chan<- app.Message) error {
	for {

		message := app.Message{}

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return err
		}

		//sync.Once()

		fmt.Printf("new message %+v\n", message)

		in <- message

	}
}

func (ws *WebSocketConnector) GetMessages() (<-chan app.Message, <-chan *app.Event, error) {

	return ws.InputMessage, ws.Events, nil
}

func (ws *WebSocketConnector) SendMessage(m app.Message, connect app.Connect) error {
	return nil
}

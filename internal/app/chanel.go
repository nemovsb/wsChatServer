package app

import (
	"sync"
	"time"
)

type Chanel struct {
	mx           *sync.Mutex
	Id           int64
	Name         string
	ClientsCount *int
	Connects     map[string]Connect
}

func NewChanel(c Connect, m Message) Chanel {
	connects := make(map[string]Connect)
	connects[m.ClientName] = c
	cc := 1
	return Chanel{
		Id:           time.Now().Unix(),
		Name:         m.ChanelName,
		ClientsCount: &cc,
		Connects:     connects,
	}
}

func (c *Chanel) AddConn(conn *Connect) {
	c.mx.Lock()

	*(c.ClientsCount)++
	c.Connects[conn.Nickname] = *conn

	c.mx.Unlock()
}

func (c *Chanel) DelConn(conn *Connect) {
	c.mx.Lock()

	*(c.ClientsCount)--
	delete(c.Connects, conn.Nickname)

	c.mx.Unlock()
}

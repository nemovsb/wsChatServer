package app

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Chanel struct {
	in        chan Message
	ctx       context.Context
	cancelCtx context.CancelFunc

	mx   *sync.Mutex
	Id   int64
	Name string

	ClientsCount *int
	Connects     map[string]Connect
}

func NewChanel(c Connect, m Message) *Chanel {
	inMessages := make(chan Message)
	connects := make(map[string]Connect)
	connects[m.ClientName] = c
	cc := 1
	ctx, cancel := context.WithCancel(context.Background())
	return &Chanel{
		in:           inMessages,
		ctx:          ctx,
		cancelCtx:    cancel,
		mx:           &sync.Mutex{},
		Id:           time.Now().Unix(),
		Name:         m.ChanelName,
		ClientsCount: &cc,
		Connects:     connects,
	}
}

func (c *Chanel) Broadcast() {
	for {

		fmt.Printf("\nready to broadcast messages! chanel: %s, %v\n", c.Name, &c.in)

		select {
		case <-c.ctx.Done():
			{

				fmt.Printf("Clients count: %d. Channel %s close.\n", *c.ClientsCount, c.Name)

				return
			}
		case message := <-c.in:
			{

				for _, connect := range c.Connects {

					if connect.Nickname == message.ClientName {
						continue
					}

					err := send(connect, message)
					if err != nil {
						if c.DelConn(&connect) {
							break
						}
						fmt.Printf("send error: %s\n", err)
						continue
					}
				}
			}
		}

	}

}

func (c *Chanel) AddConn(conn *Connect) {
	c.mx.Lock()

	*(c.ClientsCount)++
	c.Connects[conn.Nickname] = *conn

	fmt.Printf("add new connect! nicname: %s, chanelClientCount: %d\n", conn.Nickname, *c.ClientsCount)

	c.mx.Unlock()
}

func (c *Chanel) DelConn(conn *Connect) (delChannelFlag bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	*(c.ClientsCount)--
	if *c.ClientsCount == 0 {

		delChannelFlag = true
	}
	fmt.Printf("delete connect! nicname: %s, chanelClientCount: %d\n", conn.Nickname, *c.ClientsCount)

	delete(c.Connects, conn.Nickname)

	return delChannelFlag
}

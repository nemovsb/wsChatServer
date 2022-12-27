package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

func (h *Handler) wsEndpoint(ctx *gin.Context) {

	fmt.Printf("wsEndpoint start\n")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	fmt.Printf("wsEndpoint CheckOrigin\n")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	fmt.Printf("wsEndpoint upgrader.Upgrade:\n")

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	fmt.Printf("wsEndpoint h.app.Serve(*ws)\n")

	h.app.Serve(*ws)

}

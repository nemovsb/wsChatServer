package router

import (
	"ws_server/internal/app"
	"ws_server/internal/app/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	app *app.App
	wsController
}

type wsController interface {
	AddConn(wsConnect *websocket.Conn) *ws.Connect
}

func NewHandler(a *app.App, wsC wsController) *Handler {
	return &Handler{
		app:          a,
		wsController: wsC,
	}
}

func NewRouter(h *Handler) (router *gin.Engine) {
	router = gin.Default()

	router.GET("", h.wsEndpoint)

	return router
}

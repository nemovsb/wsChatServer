package router

import (
	"ws_server/internal/app"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	app *app.App
	wsController
}

type wsController interface {
	ServeConnection(wsConnect *websocket.Conn)
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

package router

import (
	"ws_server/internal/app"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	app *app.App
}

func NewHandler(a *app.App) *Handler {
	return &Handler{
		app: a,
	}
}

func NewRouter(h *Handler) (router *gin.Engine) {
	router = gin.Default()

	router.GET("", h.wsEndpoint)

	return router
}

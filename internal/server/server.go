package server

import (
	"context"
	"net/http"
)

type Server struct {
	server http.Server
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}

func NewServer(cfg ServerConfig, handler http.Handler) *Server {
	//fmt.Printf("Port : %s", fmt.Sprint(":"+conf.Port))
	serv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	return &Server{server: serv}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ws_server/internal/app"
	"ws_server/internal/config"
	"ws_server/internal/server"
	"ws_server/internal/server/router"

	group "github.com/oklog/run"
)

var ErrOsSignal = errors.New("got os signal")

func main() {

	config, err := config.ViperConfigurationProvider(os.Getenv("GOLANG_ENVIRONMENT"), false)
	if err != nil {
		log.Fatal("Read config error: ", err)
	}

	var (
		serviceGroup        group.Group
		interruptionChannel = make(chan os.Signal, 1)
	)

	servConfig := server.NewServerConfig(config.HttpServer.Port)

	application := app.NewApp()
	go func() {
		application.ServeChanels()
	}()

	handler := router.NewHandler(application)

	router := router.NewRouter(handler)

	server := server.NewServer(servConfig, router)

	serviceGroup.Add(func() error {
		signal.Notify(interruptionChannel, syscall.SIGINT, syscall.SIGTERM)
		osSignal := <-interruptionChannel

		return fmt.Errorf("%w: %s", ErrOsSignal, osSignal)
	}, func(error) {
		interruptionChannel <- syscall.SIGINT
	})

	serviceGroup.Add(func() error {
		log.Println("HTTP server started")

		return server.Run()
	}, func(error) {
		// Graceful shutdown

		application.CloseConnections()
		log.Println("Connections closed")

		err = server.Shutdown()
		log.Println("shutdown Http Server: ", err)
	})

	err = serviceGroup.Run()
	log.Println("services stopped: ", err)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ssych/file_service/pkg/rest"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := rest.NewServer(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to init server-> %v", err))
	}

	go func() {
		if err = server.Run(); err != nil {
			panic(fmt.Sprintf("failed to run server-> %v", err))
		}
	}()

	done := make(chan bool, 1)
	go func() {
		for range sigs {
			server.Close()
			cancel()
			log.Println("Shutdown server")
			done <- true
		}
	}()
	<-done
}

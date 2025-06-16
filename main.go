package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// Create a channel to listen for OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// Context for keeping track of Server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Go Routine that listens to system calls
	go func() {
		oscall := <-sig
		fmt.Printf("system call: %+v\n", oscall)
		cancel()
	}()

	// Start server (server.go will handle port 3000)
	CreateServer(ctx)
}



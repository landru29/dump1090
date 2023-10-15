package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-s
		cancel()

		os.Exit(0)
	}()

	if err := rootCommand().ExecuteContext(ctx); err != nil {
		panic(err)
	}
}

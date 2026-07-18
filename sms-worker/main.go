package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sms-worker/statrup"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := statrup.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

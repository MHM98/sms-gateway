package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	app, err := newApplication(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// the code will block here
	runErr := app.Run(ctx)
	closeErr := app.Close()
	if err := errors.Join(runErr, closeErr); err != nil {
		log.Fatal(err)
	}
}

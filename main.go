package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("vars.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	e := app.Router()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool)
	errCh := make(chan error)

	go func() {
		<-signalChan

		if app.AssistantID != "" {
			if _, err := app.Client.DeleteAssistant(context.Background(), app.AssistantID); err != nil {
				errCh <- err
			}
		}

		if app.FileID != "" {
			if err := app.Client.DeleteFile(context.Background(), app.FileID); err != nil {
				errCh <- err
			}
		}

		done <- true
	}()

	go func() {
		if err := e.Start(app.Addr); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-done:

	case err := <-errCh:
		app.Client.DeleteAssistant(context.Background(), app.AssistantID)
		app.Client.DeleteFile(context.Background(), app.FileID)
		log.Fatal(err)
	}
}

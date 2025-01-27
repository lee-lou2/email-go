package main

import (
	"email-go/api"
	"email-go/cmd/dispatcher"
	"email-go/cmd/sender"
	"email-go/config"
	"github.com/getsentry/sentry-go"
)

func main() {
	// Sentry
	_ = sentry.Init(sentry.ClientOptions{
		Dsn: config.GetEnv("SENTRY_DSN"),
	})

	// Email Consumer
	go sender.ConsumeSend()
	go sender.ConsumePostSend()

	// Message Consumer
	go dispatcher.Run()
	// HTTP Server
	api.Run()
}

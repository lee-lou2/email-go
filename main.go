package main

import (
	"aws-ses-sender-go/api"
	"aws-ses-sender-go/cmd/dispatcher"
	"aws-ses-sender-go/cmd/sender"
	"aws-ses-sender-go/config"
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

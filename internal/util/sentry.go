package util

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
)

func StartSentry(dsn string) error {
	releaseVersion := os.Getenv("SOURCE_COMMIT")
	log.Printf("Release version: %s", releaseVersion)
	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            false,
		SendDefaultPII:   true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true,
		AttachStacktrace: true,
		SampleRate:       1.0,
		Release:          releaseVersion,
	})
}

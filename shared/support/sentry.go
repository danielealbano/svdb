package shared_support

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"time"
)

type SentryWrappedFunc func()

func getSentryDsn() string {
	return os.Getenv("SENTRY_DSN")
}

func shouldEnableSentry() bool {
	return getSentryDsn() != ""
}

func setupSentry() {
	var sentryDsn string
	var hostname string
	var err error

	if !shouldEnableSentry() {
		Logger().Info().Msg("SENTRY_DSN is not set, Sentry will not be initialized.")
		return
	}

	sentryDsn = getSentryDsn()

	Logger().Info().Msgf("initializing sentry with dsn: %s", sentryDsn)

	// Read the hostname of the machine
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		TracesSampleRate: 1.0,
		Debug:            false,
		Environment:      os.Getenv("SENTRY_APP_ENVIRONMENT"),
		Release:          GetVersion(),
		ServerName:       hostname,
		Transport:        sentry.NewHTTPSyncTransport(),
	})
	if err != nil {
		Logger().Fatal().Msgf("failed to initialize Sentry: %s", err)
	}
}

func WrapWithSentry(callback SentryWrappedFunc) {
	// Setup Sentry
	setupSentry()
	defer sentry.Flush(2 * time.Second)

	// Report the panic if the program panics
	defer func() {
		var errAny interface{}
		var errInt error
		var errStr string
		var ok bool
		if errAny = recover(); errAny == nil {
			return
		}

		if errInt, ok = errAny.(error); !ok {
			if errStr, ok = errAny.(string); ok {
				errInt = errors.New(errStr)
			} else {
				errInt = errors.New("unknown panic")
			}
		}

		// Report the panic to Sentry
		sentry.CaptureException(errInt)

		// Report the panic on the console
		panic(errInt)
	}()

	callback()
}

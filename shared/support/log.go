package shared_support

import (
	"fmt"
	"github.com/jwalton/go-supportscolor"
	"github.com/phuslu/log"
	"golang.org/x/net/context"
	"io"
	"os"
	"runtime"
	"strings"
)

const (
	ansiReset          = "\x1b[0m"
	ansiFgColorRed     = "\x1b[31m"
	ansiFgColorYellow  = "\x1b[33m"
	ansiFgColorMagenta = "\x1b[35m"
	ansiFgColorCyan    = "\x1b[36m"
	ansiBold           = "\x1b[1m"
)

type logMessage struct {
	message string
	writer  io.Writer
}

var (
	alreadySetupLogger = false
	rootPath           = ""
)

func init() {
	_, currentFilePath, _, _ := runtime.Caller(0)
	rootPath = strings.TrimSuffix(strings.Split(currentFilePath, "/shared/support/log.go")[0], "/")
}

func absPathToRelativeRootPath(absPath string) string {
	return strings.TrimPrefix(strings.Replace(absPath, rootPath, "", 1), "/")
}

func loggerRoutine(ctx context.Context, logMessagesChan <-chan logMessage) {
	exitFromLoop := false
	for !exitFromLoop {
		select {
		case <-ctx.Done():
			exitFromLoop = true
			continue
		case m := <-logMessagesChan:
			_, err := fmt.Fprintln(m.writer, m.message)
			if err != nil {
				panic(err)
			}
		}
	}

	for m := range logMessagesChan {
		_, err := fmt.Fprintln(m.writer, m.message)
		if err != nil {
			panic(err)
		}
	}
}

func SetupLogger(ctx context.Context) {
	if alreadySetupLogger {
		return
	}
	alreadySetupLogger = true

	logMessagesChannel := make(chan logMessage, 1000)
	go loggerRoutine(ctx, logMessagesChannel)

	logLevel := log.InfoLevel

	if DelveEnabled {
		logLevel = log.TraceLevel
	}

	log.DefaultLogger = log.Logger{
		Level:  logLevel,
		Caller: -1,
		Writer: &log.ConsoleWriter{
			ColorOutput: true,
			Writer:      os.Stdout,
			Formatter: func(w io.Writer, args *log.FormatterArgs) (n int, err error) {
				var sb strings.Builder
				var ansiStart = ""
				var ansiEnd = ""
				var firstKeyValuePrinted = false

				// Check if stdout supports colors, it assumes that the os.Stdout is passed to the writer argument of
				// the logger ConsoleWriter writer
				if supportscolor.Stdout().SupportsColor {
					switch args.Level {
					case "trace":
						ansiStart = ansiFgColorMagenta
					case "debug":
						ansiStart = ansiFgColorCyan
					case "warn":
						ansiStart = ansiFgColorYellow
					case "error":
						ansiStart = ansiFgColorRed + ansiBold
					case "fatal":
						ansiStart = ansiFgColorRed + ansiBold
					case "panic":
						ansiStart = ansiFgColorRed + ansiBold
					}
					ansiEnd = ansiReset
				}

				// Prints out the beginning of the log line
				sb.WriteString(fmt.Sprintf(
					"%s[%s][GOROUTINE:%s][%s] ",
					ansiStart,
					args.Time,
					args.Goid,
					strings.ToUpper(args.Level)))

				// If there is a caller defined try to print it out
				if args.Caller != "" {
					sb.WriteString(fmt.Sprintf("%s > ", absPathToRelativeRootPath(args.Caller)))
				}

				// Print out the message
				sb.WriteString(args.Message)

				// Print the additional data
				for _, kv := range args.KeyValues {
					if kv.Key == "callerfunc" {
						continue
					}

					if !firstKeyValuePrinted {
						sb.WriteString(" (")
					} else {
						sb.WriteString(", ")
					}

					sb.WriteString(fmt.Sprintf("%s=%v", kv.Key, kv.Value))

					firstKeyValuePrinted = true
				}

				if firstKeyValuePrinted {
					sb.WriteString(")")
				}

				// If there is a stack trace print it out
				if args.Stack != "" {
					sb.WriteString("\n")
					sb.WriteString(args.Stack)
				}

				// Print the ansi reset
				sb.WriteString(ansiEnd)

				// Push the message to the log messages channel
				logMessagesChannel <- logMessage{
					message: sb.String(),
					writer:  w,
				}

				return sb.Len(), nil
			},
		},
	}
}

func Logger() *log.Logger {
	return &log.DefaultLogger
}

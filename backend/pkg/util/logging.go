package util

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lokiClient LokiClient

// Parse CMD flags and configure logger accordingly
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debug := flag.Bool("debug", false, "sets log level to debug")
	rawLogs := flag.Bool("pretty-logs", false, "format logs for human readabiliy rather than parsing")

	flag.Parse()
	if *rawLogs {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}

		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}

	// Default level is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msg("Starting with loglevel debug enabled")

	// Send logs every 10 seconds or if 500 log events have been collected
	lokiClient = LokiClient{
		PushIntveralSeconds: 10,
		MaxBatchSize:        500,
		LokiEndpoint:        GetStringEnvWithDefault("LOKI_HOST", "http://localhost:3100"),
	}

	go lokiClient.bgRun()

	// Include calling line in code in log
	log.Logger = log.With().Caller().Logger().Hook(LokiHook{})
}

type LokiHook struct {
}

func (h LokiHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	lokiClient.Values = append(lokiClient.Values, []string{strconv.FormatInt(time.Now().UnixNano(), 10), msg})
}

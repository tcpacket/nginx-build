package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	log   *zerolog.Logger
	setup = &sync.Once{}
)

func Get() *zerolog.Logger {
	setup.Do(func() {
		l := zerolog.NewConsoleWriter()
		l.Out = os.Stdout
		l.TimeFormat = time.RFC3339
		lg := zerolog.New(l).With().Timestamp().Logger()
		log = &lg
	})
	return log
}

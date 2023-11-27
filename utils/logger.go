package utils

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

var (
	loggerOnce = sync.Once{}
	_Logger    zerolog.Logger
)

func GetLogger() *zerolog.Logger {
	loggerOnce.Do(factoryLogger)
	return &_Logger
}

func factoryLogger() {
	// TODO Should be switched by prod/dev env
	out := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("[messenger] %s", i)
		},
	}
	zerologLvl := zerolog.Level(GetConfig().LogLevel)
	_Logger = zerolog.New(out).Level(zerologLvl).With().Timestamp().Logger()
}

package utils

import "github.com/rs/zerolog/log"

type IntegratedLogger struct {
}

func (i IntegratedLogger) Printf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v)
}

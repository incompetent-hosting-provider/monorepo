package helper

import (
	"os"

	"github.com/rs/zerolog/log"
)

type Error interface {
	Error() string
}

type CustomError struct {
	Msg string
}

func NewCustomError(msg string) CustomError {
	return CustomError{Msg: msg}
}

func (e CustomError) Error() string {
	return e.Msg
}

func HandleFatalError(e error, msg string) {
	if e != nil {
		log.Fatal().Msgf("%s: %s", e, msg)
		os.Exit(1)
	}
}

func HandleError(e error, msg string) {
	if e != nil {
		log.Error().Msgf("%s: %s", e, msg)
	}
}

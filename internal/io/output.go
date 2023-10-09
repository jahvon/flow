package io

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

type OutputFormat string

const (
	OutputFormatJSON       OutputFormat = "json"
	OutputFormatPrettyJson OutputFormat = "jsonp"
	OutputFormatYAML       OutputFormat = "yaml"
	OutputFormatDefault    OutputFormat = "default"
)

func Log() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
	context := zerolog.New(writer).With().Timestamp()
	return context.Logger()
}

func PrintInfo(message string) {
	log := Log()
	log.Info().Msg(message)
}

func PrintNotice(message string) {
	color.HiBlue(message)
}

func PrintQuestion(question string) {
	color.HiCyan(question)
}

func PrintSuccess(message string) {
	color.HiGreen(message)
}

func PrintWarning(message string) {
	color.HiYellow(message)
}

func PrintErrorAndExit(err error) {
	PrintError(err)
	os.Exit(1)
}

func PrintError(err error) {
	color.HiRed(err.Error())
}

type StdOutWriter struct {
	LogAsDebug bool
}

func (w StdOutWriter) Write(p []byte) (n int, err error) {
	log := Log()
	trimmedP := strings.TrimSpace(string(p))
	if trimmedP == "" {
		return len(p), nil
	}

	if w.LogAsDebug {
		log.Debug().Msg(trimmedP)
		return len(p), nil
	}

	log.Info().Msg(trimmedP)
	return len(p), nil
}

type StdErrWriter struct {
	LogAsDebug bool
}

func (w StdErrWriter) Write(p []byte) (n int, err error) {
	log := Log()
	trimmedP := strings.TrimSpace(string(p))
	if trimmedP == "" {
		return len(p), nil
	}

	if w.LogAsDebug {
		log.Debug().Msg(string(p))
		return len(p), nil
	}

	log.Error().Msg(string(p))
	return len(p), nil
}

package io

import (
	"os"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

func Log() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
	context := zerolog.New(writer).With().Timestamp()
	return context.Logger()
}

func PrintInfo(message string) {
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
